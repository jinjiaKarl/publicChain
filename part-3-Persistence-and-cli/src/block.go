package src

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// 将 Block 序列化为一个字节数组
func (b *Block) Serialize() []byte {
	//定义一个 buffer 存储序列化之后的数据
	var result bytes.Buffer
	// 生成一个Encoder
	encoder := gob.NewEncoder(&result)
	//对 block 进行编码
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	//返回字节数组
	return result.Bytes()
}

// 将字节数组反序列化为一个 Block
func DeserializeBlock(d []byte) *Block {
	//定义一个block 存储反序列化之后的数据
	var block Block
	//生成一个解码器Decoder
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

//生成一个新的区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}
