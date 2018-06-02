package src

import (
	"github.com/bolt"
	"log"
	"fmt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

/*
 *1. 打开一个数据库文件
 *2. 检查文件里面是否已经存储了一个区块链
 *3. 如果已经存储了一个区块链：
    1. 创建一个新的 `Blockchain` 实例
    2. 设置 `Blockchain` 实例的 tip 为数据库中存储的最后一个块的哈希
 *4. 如果没有区块链：
    1. 创建创世块
    2. 存储到数据库
    3. 将创世块哈希保存为最后一个块的哈希
    4. 创建一个新的 `Blockchain` 实例，初始时 tip 指向创世块（tip 有尾部，尖端的意思，在这里 tip 存储的是最后一个块的哈希）
 */
// tip 这个词本身有事物尖端或尾部的意思，这里指的是存储最后一个块的哈希
// 在链的末端可能出现短暂分叉的情况，所以选择 tip 其实也就是选择了哪条链
// db 存储数据库连接
type Blockchain struct {
	Tip []byte
	Db  *bolt.DB
}

func NewBlockchain() *Blockchain {
	var tip []byte
	// 打开一个 BoltDB 文件
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//在 BoltDB 中，数据库操作通过一个事务（transaction）进行操作。
	// 有两种类型的事务：只读（read-only）和读写（read-write）。
	// 这里，打开的是一个读写事务（`db.Update(...)`），因为我们可能会向数据库中添加创世块。
	err = db.Update(func(tx *bolt.Tx) error {
		//在这里，我们先获取了存储区块的 bucket：如果存在，就从中读取 `l` 键；
		// 如果不存在，就生成创世块，创建 bucket，并将区块保存到里面，然后更新 `l` 键以存储链中最后一个块的哈希。
		b := tx.Bucket([]byte(blocksBucket))

		// 如果数据库中不存在区块链就创建一个，否则直接读取最后一个块的哈希
		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// 加入区块时，需要将区块持久化到数据库中
func (bc *Blockchain) AddBlock(data string) {

	var lastHash []byte
	//这是 BoltDB 事务的另一个类型（只读）
	// 首先获取最后一个块的哈希用于生成新块的哈希
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)
	//更新bucket中"l"的键值
	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.Tip = newBlock.Hash

		return nil
	})
}

// BoltDB 允许对一个 bucket 里面的所有 key 进行迭代，但是所有的 key 都以字节序进行存储，
// 而且我们想要以区块能够进入区块链中的顺序进行打印。
// 此外，因为我们不想将所有的块都加载到内存中（因为我们的区块链数据库可能很大！或者现在可以假装它可能很大），
// 我们将会一个一个地读取它们。故而，我们需要一个区块链迭代器（`BlockchainIterator`）
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.Tip, bc.Db}

	return bci
}

// 返回链中的下一个块
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}
