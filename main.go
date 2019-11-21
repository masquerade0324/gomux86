package main

import (
	"fmt"
	"io"
	"os"
)

/* memory size: 1MB */
const memSize = 1024 * 1024

const (
	regEax = iota
	regEcx
	regEdx
	regEbx
	regEsp
	regEbp
	regEsi
	regEdi
	regCnt
)

var regNames = [regCnt]string{"EAX", "ECX", "EDX", "EBX", "ESP", "EBP", "ESI", "EDI"}

type Emulator struct {
	registers [regCnt]uint32
	eflags    uint32
	memory    []uint8
	eip       uint32
}

func NewEmulator(size int, eip, esp uint32) *Emulator {
	emu := new(Emulator)
	emu.memory = make([]uint8, size)
	emu.eip = eip
	emu.registers[regEsp] = esp
	return emu
}

func dumpRegisters(emu *Emulator) {
	for i, v := range emu.registers {
		fmt.Printf("%v = %08x\n", regNames[i], v)
	}
	fmt.Printf("EIP = %08x\n", emu.eip)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: gomux86 filename")
		os.Exit(1)
	}

	emu := NewEmulator(memSize, 0x00007C00, 0x00007C00)

	binary, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Can't open this file: %v", os.Args[1])
		os.Exit(1)
	}
	defer binary.Close()

	reader := io.LimitReader(binary, 512)
	if _, err := reader.Read(emu.memory[0x00007c00:]); err != nil {
		os.Exit(1)
	}

	initInstructions()

	for emu.eip < memSize {
		code := getCodeU8(emu, 0)

		fmt.Printf("EIP = %x, Code = %02x\n", emu.eip, code)

		if instructions[code] == nil {
			fmt.Printf("\n\nNot Implemented: %x\n", code)
			break
		}

		/* execute the instruction */
		instructions[code](emu)

		if emu.eip == 0x00 {
			fmt.Printf("\n\nend of program.\n\n")
			break
		}
	}

	dumpRegisters(emu)
}
