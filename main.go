package main

import (
	"fmt"
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

var instructions [256]func(*Emulator)

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

func getCodeU8(emu *Emulator, index int) uint8 {
	return emu.memory[int(emu.eip)+index]
}

func getCodeS8(emu *Emulator, index int) int8 {
	return int8(emu.memory[int(emu.eip)+index])
}

func getCodeU32(emu *Emulator, index int) uint32 {
	var ret uint32 = 0
	for i := 0; i < 4; i++ {
		ret += uint32(getCodeU8(emu, index+i)) << uint((i * 8))
	}
	return ret
}

func movR32Imm32(emu *Emulator) {
	reg := getCodeU8(emu, 0) - 0xB8
	val := getCodeU32(emu, 1)
	emu.registers[reg] = val
	emu.eip += 5
}

func shortJmp(emu *Emulator) {
	diff := getCodeS8(emu, 1)
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func initInstructions() {
	for i := 0; i < 8; i++ {
		instructions[0xB8+i] = movR32Imm32
	}
	instructions[0xEB] = shortJmp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: gomux86 filename")
		os.Exit(1)
	}

	emu := NewEmulator(memSize, 0x00000000, 0x00007C00)

	binary, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Can't open this file: %v", os.Args[1])
		os.Exit(1)
	}
	defer binary.Close()

	bs := make([]uint8, 512)
	if _, err := binary.Read(bs); err != nil {
		os.Exit(1)
	}
	emu.memory = bs

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
