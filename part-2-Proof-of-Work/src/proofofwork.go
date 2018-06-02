package src

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

// 难度值，这里表示哈希的前 24 位必须是 0
const targetBits = 16

const maxNonce = math.MaxInt64

// 每个块的工作量都必须要证明，所有有个指向 Block 的指针
// target 是目标，我们最终要找的哈希必须要小于目标
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// target 等于 1 左移 256 - targetBits 位
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	//左移 256 - targetBits 位
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

// 工作量证明用到的数据有：PrevBlockHash, Data, Timestamp, targetBits, nonce
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			//将一个 int64 转化为一个字节数组(byte array)
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		//多个字节数组间的连接符
		[]byte{},
	)

	return data
}

// 工作量证明的核心就是寻找有效哈希
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	//256位，32字节
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		//1.准备数据
		data := pow.prepareData(nonce)
		//2.对数据就进行哈希运算
		hash = sha256.Sum256(data)
		//3.将哈希转换成一个大整数
		hashInt.SetBytes(hash[:])
		//将这个大整数与目标进行比较，若小，则挖矿成功
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// 验证工作量，只要哈希小于目标就是有效工作量
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
