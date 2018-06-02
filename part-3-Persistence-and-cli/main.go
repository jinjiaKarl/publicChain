package main

import "publicChain/part-3-Persistence-and-cli/src"

func main() {
	//注意，无论提供什么命令行参数，都会创建一个新的链。如何去修改这个问题？
	bc := src.NewBlockchain()
	defer bc.Db.Close()

	cli := src.CLI{bc}
	cli.Run()
}
