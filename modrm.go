package main

import (
	"fmt"
	"os"
)

type modRM struct {
	mod    uint8
	reg    uint8
	rm     uint8
	sib    uint8
	disp8  int8
	disp32 uint32
}

func parseModRM(emu *Emulator, modrm *modRM) {
	code := getCodeU8(emu, 0)

	modrm.mod = (code & 0xC0) >> 6
	modrm.reg = (code & 0x38) >> 3
	modrm.rm = code & 0x07

	emu.eip += 1

	if modrm.mod != 3 && modrm.rm == 4 {
		modrm.sib = getCodeU8(emu, 0)
		emu.eip += 1
	}

	if modrm.mod == 0 && modrm.rm == 5 || modrm.mod == 2 {
		modrm.disp32 = getCodeU32(emu, 0)
		emu.eip += 4
	} else if modrm.mod == 1 {
		modrm.disp8 = getCodeS8(emu, 0)
		emu.eip += 1
	}
}

func setR8(emu *Emulator, modrm *modRM, val uint8) {
	setRegister8(emu, int(modrm.reg), val)
}

func setR32(emu *Emulator, modrm *modRM, val uint32) {
	setRegister32(emu, int(modrm.reg), val)
}

func setRm8(emu *Emulator, modrm *modRM, val uint8) {
	if modrm.mod == 3 {
		setRegister8(emu, int(modrm.rm), val)
	} else {
		addr := calcMemoryAddr(emu, modrm)
		setMemory8(emu, addr, uint32(val))
	}
}

func setRm32(emu *Emulator, modrm *modRM, val uint32) {
	if modrm.mod == 3 {
		setRegister32(emu, int(modrm.rm), val)
	} else {
		addr := calcMemoryAddr(emu, modrm)
		setMemory32(emu, addr, val)
	}
}

func getR8(emu *Emulator, modrm *modRM) uint8 {
	return getRegister8(emu, int(modrm.reg))
}

func getR32(emu *Emulator, modrm *modRM) uint32 {
	return getRegister32(emu, int(modrm.reg))
}

func getRm8(emu *Emulator, modrm *modRM) uint8 {
	if modrm.mod == 3 {
		return getRegister8(emu, int(modrm.rm))
	} else {
		addr := calcMemoryAddr(emu, modrm)
		return uint8(getMemory8(emu, addr))
	}
}

func getRm32(emu *Emulator, modrm *modRM) uint32 {
	if modrm.mod == 3 {
		return getRegister32(emu, int(modrm.rm))
	} else {
		addr := calcMemoryAddr(emu, modrm)
		return getMemory32(emu, addr)
	}
}

func calcMemoryAddr(emu *Emulator, modrm *modRM) uint32 {
	var addr uint32
	if modrm.mod == 0 {
		if modrm.rm == 4 {
			fmt.Printf("not implemented ModRM mod = 0, rm = 4\n")
			os.Exit(1)
		} else if modrm.rm == 5 {
			addr = modrm.disp32
		} else {
			addr = getRegister32(emu, int(modrm.rm))
		}
	} else if modrm.mod == 1 {
		if modrm.rm == 4 {
			fmt.Printf("not implemented ModRM mod = 1, rm = 4\n")
			os.Exit(1)
		} else {
			addr = uint32(int32(getRegister32(emu, int(modrm.rm))) + int32(modrm.disp8))
		}
	} else if modrm.mod == 2 {
		if modrm.rm == 4 {
			fmt.Printf("not implemented ModRM mod = 2, rm = 4\n")
			os.Exit(1)
		} else {
			addr = getRegister32(emu, int(modrm.rm)) + modrm.disp32
		}
	} else {
		fmt.Printf("not implemented ModRM mod = 3\n")
		os.Exit(1)
	}
	return addr
}
