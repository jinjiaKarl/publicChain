package src

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	Bc *Blockchain
}

const usage = `
Usage:
  addblock -data BLOCK_DATA    add a block to the blockchain
  printchain                   print all the blocks of the blockchain
`

func (cli *CLI) printUsage() {
	fmt.Println(usage)
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()
	//使用flag包来解析命令行参数
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	//给 `addblock` 添加 `-data` 标志。`printchain` 没有任何标志。
	addBlockData := addBlockCmd.String("data", "", "Block data")
	//os.Args用于获取命令行参数
	switch os.Args[1] {
	case "addblock":
		//调用 flag.Parse() 解析命令行参数到定义的 flag
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	//我们检查用户提供的命令，解析相关的 `flag` 子命令
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.Bc.AddBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
