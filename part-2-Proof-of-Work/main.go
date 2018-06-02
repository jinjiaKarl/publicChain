package main

import (
	"fmt"
	"strconv"
	"publicChain/part-2-Proof-of-Work/src"
)

func main() {

	bc := src.NewBlockChain()
	bc.AddBlock("Send 1 BTC to jinjia")
	bc.AddBlock("Send 2 more BTC to jinjia")

	//遍历区块
	for _, block := range bc.Blocks {
		fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		//验证工作量
		pow := src.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))

		fmt.Println()
	}
}
