package main

import "fmt"

var biosToTerminal = [8]int{30, 34, 32, 36, 31, 35, 33, 37}

func putString(s string, n int) {
	for i := 0; i < n; i++ {
		ioOut8(0x03f8, s[i])
	}
}

func biosVideoTeletype(emu *Emulator) {
	color := getRegister8(emu, regBl) & 0x0f
	char := getRegister8(emu, regAl)

	terminalColor := biosToTerminal[color&0x07]
	bright := 0
	if color&0x08 != 0 {
		bright = 1
	}
	buf := fmt.Sprintf("\x1b[%d;%dm%c\x1b[m", bright, terminalColor, char)
	putString(buf, len(buf))
}

func biosVideo(emu *Emulator) {
	fun := getRegister8(emu, regAh)
	if fun == 0x0e {
		biosVideoTeletype(emu)
	} else {
		fmt.Printf("not implemented BIOS vidoe function: 0x%02x\n", fun)
	}
}
