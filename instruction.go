package main

import (
	"fmt"
	"os"
)

var instructions [256]func(*Emulator)

func movR8Imm8(emu *Emulator) {
	reg := getCodeU8(emu, 0) - 0xB0
	setRegister8(emu, int(reg), getCodeU8(emu, 1))
	emu.eip += 2
}

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

func movR8Rm8(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)
	rm8 := getRm8(emu, modrm)
	setR8(emu, modrm, rm8)
}

func movR32Rm32(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)
	rm32 := getRm32(emu, modrm)
	setR32(emu, modrm, rm32)
}

func movRm8R8(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)
	r8 := getR8(emu, modrm)
	setRm8(emu, modrm, r8)
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

func cmpR32Rm32(emu *Emulator) {
	emu.eip += 1
	modrm := new(modRM)
	parseModRM(emu, modrm)
	r32 := getR32(emu, modrm)
	rm32 := getRm32(emu, modrm)
	res := uint64(r32) - uint64(rm32)
	updateEflagsSub(emu, r32, rm32, res)
}

func cmpRm32Imm8(emu *Emulator, modrm *modRM) {
	rm32 := getRm32(emu, modrm)
	imm8 := uint32(int32(getCodeS8(emu, 0)))
	emu.eip += 1
	res := uint64(rm32) - uint64(imm8)
	updateEflagsSub(emu, rm32, imm8, res)
}

func cmpAlImm8(emu *Emulator) {
	imm8 := getCodeU8(emu, 1)
	al := getRegister8(emu, regAl)
	res := uint64(al) - uint64(imm8)
	updateEflagsSub(emu, uint32(al), uint32(imm8), res)
	emu.eip += 2
}

func cmpEaxImm32(emu *Emulator) {
	imm32 := getCodeU32(emu, 1)
	eax := getRegister32(emu, regEax)
	res := uint64(eax) - uint64(imm32)
	updateEflagsSub(emu, eax, imm32, res)
	emu.eip += 5
}

func subRm32Imm8(emu *Emulator, modrm *modRM) {
	rm32 := getRm32(emu, modrm)
	imm8 := uint32(int32(getCodeS8(emu, 0)))
	emu.eip += 1
	res := uint64(rm32) - uint64(imm8)
	setRm32(emu, modrm, uint32(res))
	updateEflagsSub(emu, rm32, imm8, res)
}

func incR32(emu *Emulator) {
	reg := int(getCodeU8(emu, 0) - 0x40)
	setRegister32(emu, reg, getRegister32(emu, reg)+1)
	emu.eip += 1
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
	case 7:
		cmpRm32Imm8(emu, modrm)
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

func jc(emu *Emulator) {
	diff := int8(0)
	if isCarry(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func jnc(emu *Emulator) {
	diff := int8(0)
	if !isCarry(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func jz(emu *Emulator) {
	diff := int8(0)
	if isZero(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func jnz(emu *Emulator) {
	diff := int8(0)
	if !isZero(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func js(emu *Emulator) {
	diff := int8(0)
	if isSign(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func jns(emu *Emulator) {
	diff := int8(0)
	if !isSign(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func jo(emu *Emulator) {
	diff := int8(0)
	if isOverflow(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func jno(emu *Emulator) {
	diff := int8(0)
	if !isOverflow(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func jl(emu *Emulator) {
	diff := int8(0)
	if isSign(emu) != isOverflow(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
}

func jle(emu *Emulator) {
	diff := int8(0)
	if isZero(emu) || isSign(emu) != isOverflow(emu) {
		diff = getCodeS8(emu, 1)
	}
	emu.eip = uint32(int32(emu.eip) + int32(diff+2))
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

func inAlDx(emu *Emulator) {
	addr := uint16(getRegister32(emu, regEdx) & 0xffff)
	val := ioIn8(addr)
	setRegister8(emu, regAl, val)
	emu.eip += 1
}

func outDxAl(emu *Emulator) {
	addr := uint16(getRegister32(emu, regEdx) & 0xffff)
	val := getRegister8(emu, regAl)
	ioOut8(addr, val)
	emu.eip += 1
}

func swi(emu *Emulator) {
	intIndex := getCodeU8(emu, 1)
	emu.eip += 2

	if intIndex == 0x10 {
		biosVideo(emu)
	} else {
		fmt.Printf("unknown interrupt: 0x%02x\n", intIndex)
	}
}

func initInstructions() {
	instructions[0x01] = addRm32R32
	instructions[0x3B] = cmpR32Rm32
	instructions[0x3C] = cmpAlImm8
	instructions[0x3D] = cmpEaxImm32
	for i := 0; i < 8; i++ {
		instructions[0x40+i] = incR32
	}
	for i := 0; i < 8; i++ {
		instructions[0x50+i] = pushR32
	}
	for i := 0; i < 8; i++ {
		instructions[0x58+i] = popR32
	}
	instructions[0x68] = pushImm32
	instructions[0x6A] = pushImm8
	instructions[0x70] = jo
	instructions[0x71] = jno
	instructions[0x72] = jc
	instructions[0x73] = jnc
	instructions[0x74] = jz
	instructions[0x75] = jnz
	instructions[0x78] = js
	instructions[0x79] = jns
	instructions[0x7C] = jl
	instructions[0x7E] = jle
	instructions[0x83] = code83
	instructions[0x88] = movRm8R8
	instructions[0x89] = movRm32R32
	instructions[0x8A] = movR8Rm8
	instructions[0x8B] = movR32Rm32
	for i := 0; i < 8; i++ {
		instructions[0xB0+i] = movR8Imm8
	}
	for i := 0; i < 8; i++ {
		instructions[0xB8+i] = movR32Imm32
	}
	instructions[0xC3] = ret
	instructions[0xC7] = movRm32Imm32
	instructions[0xC9] = leave
	instructions[0xCD] = swi
	instructions[0xE8] = callRel32
	instructions[0xE9] = nearJmp
	instructions[0xEB] = shortJmp
	instructions[0xEC] = inAlDx
	instructions[0xEE] = outDxAl
	instructions[0xFF] = codeFF
}
