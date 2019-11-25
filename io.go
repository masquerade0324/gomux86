package main

import (
	"bufio"
	"os"
)

func ioIn8(addr uint16) uint8 {
	var ch uint8
	if addr == 0x03f8 {
		ch, _ = bufio.NewReader(os.Stdin).ReadByte()
	} else {
		ch = 0
	}
	return ch
}

func ioOut8(addr uint16, val uint8) {
	if addr == 0x03f8 {
		w := bufio.NewWriter(os.Stdout)
		w.WriteByte(val)
		w.Flush()
	}
}
