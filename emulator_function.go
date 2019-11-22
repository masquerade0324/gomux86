package main

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

func getCodeS32(emu *Emulator, index int) int32 {
	return int32(getCodeU32(emu, index))
}

func setMemory8(emu *Emulator, addr, val uint32) {
	emu.memory[addr] = uint8(val) & 0xFF
}

func setMemory32(emu *Emulator, addr, val uint32) {
	i := uint32(0)
	for ; i < 4; i++ {
		setMemory8(emu, addr+i, val>>(i*8))
	}
}

func getMemory8(emu *Emulator, addr uint32) uint32 {
	return uint32(emu.memory[addr])
}

func getMemory32(emu *Emulator, addr uint32) uint32 {
	i := uint32(0)
	ret := uint32(0)
	for ; i < 4; i++ {
		ret |= getMemory8(emu, addr+i) << (8 * i)
	}
	return ret
}

func setRegister32(emu *Emulator, index int, val uint32) {
	emu.registers[index] = val
}

func getRegister32(emu *Emulator, index int) uint32 {
	return emu.registers[index]
}

func push32(emu *Emulator, val uint32) {
	addr := getRegister32(emu, regEsp) - 4
	setRegister32(emu, regEsp, addr)
	setMemory32(emu, addr, val)
}

func pop32(emu *Emulator) uint32 {
	addr := getRegister32(emu, regEsp)
	ret := getMemory32(emu, addr)
	setRegister32(emu, regEsp, addr+4)
	return ret
}
