package src

import (
	"bytes"
	"encoding/binary"
	"log"
)

// 将一个 int64 转化为一个字节数组(byte array)
func IntToHex(num int64) []byte {

	buff := new(bytes.Buffer)
	//将num变为二进制写入buff中
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
