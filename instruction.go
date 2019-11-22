package main

import (
	"fmt"
	"os"
)

var instructions [256]func(*Emulator)

func movR32Imm32(emu *Emulator) {
	reg := getCodeU8(emu, 0) - 0xB8
	val := getCodeU32(emu, 1)
	emu.registers[reg] = val
	emu.eip += 5
}

func movRm32Imm32(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)
	val := getCodeU32(emu, 0)
	emu.eip += 4
	setRm32(emu, modrm, val)
}

func movR32Rm32(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)
	rm32 := getRm32(emu, modrm)
	setR32(emu, modrm, rm32)
}

func movRm32R32(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)
	r32 := getR32(emu, modrm)
	setRm32(emu, modrm, r32)
}

func addRm32R32(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)
	r32 := getR32(emu, modrm)
	rm32 := getRm32(emu, modrm)
	setRm32(emu, modrm, rm32+r32)
}

func addRm32Imm8(emu *Emulator, modrm *modRM) {
	rm32 := getRm32(emu, modrm)
	imm8 := uint32(int32(getCodeS8(emu, 0)))
	emu.eip += 1
	setRm32(emu, modrm, rm32+imm8)
}

func subRm32Imm8(emu *Emulator, modrm *modRM) {
	rm32 := getRm32(emu, modrm)
	imm8 := uint32(int32(getCodeS8(emu, 0)))
	emu.eip += 1
	setRm32(emu, modrm, rm32-imm8)
}

func incRm32(emu *Emulator, modrm *modRM) {
	val := getRm32(emu, modrm)
	setRm32(emu, modrm, val+1)
}

func code83(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)

	switch modrm.reg {
	case 0:
		addRm32Imm8(emu, modrm)
	case 5:
		subRm32Imm8(emu, modrm)
	default:
		fmt.Printf("not implemented: 83 /%v\n", modrm.reg)
		os.Exit(1)
	}
}

func codeFF(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)

	switch modrm.reg {
	case 0:
		incRm32(emu, modrm)
	default:
		fmt.Printf("not implemented: FF /%v\n", modrm.reg)
		os.Exit(1)
	}
}

func shortJmp(emu *Emulator) {
	diff := getCodeS8(emu, 1)
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func nearJmp(emu *Emulator) {
	diff := getCodeS32(emu, 1)
	emu.eip = uint32(int32(emu.eip) + diff + 5)
}

func pushR32(emu *Emulator) {
	reg := getCodeU8(emu, 0) - 0x50
	push32(emu, getRegister32(emu, int(reg)))
	emu.eip += 1
}

func pushImm8(emu *Emulator) {
	val := getCodeU8(emu, 1)
	push32(emu, uint32(val))
	emu.eip += 2
}

func pushImm32(emu *Emulator) {
	val := getCodeU32(emu, 1)
	push32(emu, val)
	emu.eip += 5
}

func popR32(emu *Emulator) {
	reg := getCodeU8(emu, 0) - 0x58
	setRegister32(emu, int(reg), pop32(emu))
	emu.eip += 1
}

func callRel32(emu *Emulator) {
	diff := getCodeS32(emu, 1)
	push32(emu, emu.eip+5)
	emu.eip = uint32(int32(emu.eip) + diff + 5)
}

func ret(emu *Emulator) {
	emu.eip = pop32(emu)
}

func leave(emu *Emulator) {
	ebp := getRegister32(emu, regEbp)
	setRegister32(emu, regEsp, ebp)
	setRegister32(emu, regEbp, pop32(emu))
	emu.eip += 1
}

func initInstructions() {
	instructions[0x01] = addRm32R32
	for i := 0; i < 8; i++ {
		instructions[0x50+i] = pushR32
	}
	for i := 0; i < 8; i++ {
		instructions[0x58+i] = popR32
	}
	instructions[0x68] = pushImm32
	instructions[0x6A] = pushImm8
	instructions[0x83] = code83
	instructions[0x89] = movRm32R32
	instructions[0x8B] = movR32Rm32
	for i := 0; i < 8; i++ {
		instructions[0xB8+i] = movR32Imm32
	}
	instructions[0xC3] = ret
	instructions[0xC7] = movRm32Imm32
	instructions[0xC9] = leave
	instructions[0xE8] = callRel32
	instructions[0xE9] = nearJmp
	instructions[0xEB] = shortJmp
	instructions[0xFF] = codeFF
}
