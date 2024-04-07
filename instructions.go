package main

import (
	"fmt"
	"log"
	"math/rand"
)

// 0nnn - SYS addr
// Jump to a machine code routine at nnn.
//
// This instruction is only used on the old computers on which Chip-8
// was originally implemented. It is ignored by modern interpreters.
func (c *chip8) sysAddr(addr uint16) {
	c.print("SYS %04X", addr)
	c.pc = addr
}

// 00E0 - CLS
// Clear the display.
func (c *chip8) cls() {
	c.print("CLS")
	for x := range c.screen {
		for y := range c.screen[x] {
			c.screen[x][y] = false
		}
	}
	c.pc += 2
}

// 00EE - RET
// Return from a subroutine.
//
// The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
func (c *chip8) ret() {
	c.print("RET")
	c.pc = c.stack[c.sp]
	c.sp--
	c.pc += 2
}

func (c *chip8) print(format string, args ...any) {
	infos := []struct {
		fmt   string
		value any
	}{
		{"PC: %04X", c.pc},
		{"I: %04X", c.i},
		{"DT: %02X", c.dt},
		{"ST: %02X\n", c.st},
		{"V0: %02X", c.v[0x0]},
		{"V1: %02X", c.v[0x1]},
		{"V2: %02X", c.v[0x2]},
		{"V3: %02X", c.v[0x3]},
		{"V4: %02X", c.v[0x4]},
		{"V5: %02X", c.v[0x5]},
		{"V6: %02X", c.v[0x6]},
		{"V7: %02X\n", c.v[0x7]},
		{"V8: %02X", c.v[0x8]},
		{"V9: %02X", c.v[0x9]},
		{"VA: %02X", c.v[0xA]},
		{"VB: %02X", c.v[0xB]},
		{"VC: %02X", c.v[0xC]},
		{"VD: %02X", c.v[0xD]},
		{"VE: %02X", c.v[0xE]},
		{"VF: %02X\n", c.v[0xF]},
	}

	format2 := "\n"
	args2 := []any{}

	for i, info := range infos {
		if i != 0 {
			format2 += " | "
		}
		format2 += info.fmt
		args2 = append(args2, info.value)
	}

	state := fmt.Sprintf(
		format2,
		args2...,
	)
	log.Print(fmt.Sprintf("\033[1;33m"+format+"\033[0m", args...) + state)
}

// 1nnn - JP addr
// Jump to location nnn.
//
// The interpreter sets the program counter to nnn.
func (c *chip8) jpAddr(addr uint16) {
	c.print("JP [%04X]", addr)
	c.pc = addr
}

// 2nnn - CALL addr
// Call subroutine at nnn.
//
// The interpreter increments the stack pointer, then puts the current PC
// on the top of the stack. The PC is then set to nnn.
func (c *chip8) callAddr(addr uint16) {
	c.print("CALL %04X", addr)
	c.sp++
	c.stack[c.sp] = c.pc
	c.pc = addr
}

// 3xkk - SE Vx, byte
// Skip next instruction if Vx = kk.
//
// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
func (c *chip8) seVxB(x, b uint8) {
	c.print("SE V%01X, %d", x, b)
	if x == b {
		c.pc += 2
	}
	c.pc += 2
}

// 4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
//
// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
func (c *chip8) sneVxB(x, b uint8) {
	c.print("SNE V%01X, %d", x, b)
	if x != b {
		c.pc += 2
	}
	c.pc += 2
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
//
// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
func (c *chip8) seVxVy(x, y uint8) {
	c.print("SE V%01X, V%01X", x, y)
	if c.v[x] == c.v[y] {
		c.pc += 2
	}
	c.pc += 2
}

// 6xkk - LD Vx, byte
// Set Vx = kk.
//
// The interpreter puts the value kk into register Vx.
func (c *chip8) ldVxB(x, b uint8) {
	c.print("LD V%01X, %d", x, b)
	c.v[x] = b
	c.pc += 2
}

//
//

// 7xkk - ADD Vx, byte
// Set Vx = Vx + kk.
//
// Adds the value kk to the value of register Vx, then stores the result in Vx.
func (c *chip8) addVxB(x, b uint8) {
	c.print("ADD V%01X, %d", x, b)
	c.v[x] += b
	c.pc += 2
}

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
//
// Stores the value of register Vy in register Vx.
func (c *chip8) ldReg(x, y uint8) {
	c.v[x] += c.v[y]
	c.pc += 2
}

// 8xy1 - OR Vx, Vy
// Set Vx = Vx OR Vy.
//
// Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx.
// A bitwise OR compares the corrseponding bits from two values, and if either bit
// is 1, then the same bit in the result is also 1. Otherwise, it is 0.
func (c *chip8) orVxVy(x, y uint8) {
	c.print("OR V%01X, V%01X", x, y)
	c.v[x] |= c.v[y]
	c.pc += 2
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
//
// Performs a bitwise AND on the values of Vx andVxVy Vy, then stores the result in Vx. A bitwise AND compares the corrseponding bits from two values, andVxVy if both bits are 1, then the same bit in the result is also 1. Otherwise, it is 0.
func (c *chip8) andVxVy(x, y uint8) {
	c.print("AND V%01X, V%01X", x, y)
	c.v[x] &= c.v[y]
	c.pc += 2
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
//
// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx. An exclusive OR compares the corrseponding bits from two values, and if the bits are not both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
func (c *chip8) xorVxVy(x, y uint8) {
	c.print("XOR V%01X, V%01X", x, y)
	c.v[x] ^= c.v[y]
	c.pc += 2
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
//
// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
func (c *chip8) addVxVy(x, y uint8) {
	c.print("ADD V%01X, V%01X", x, y)
	vx := c.v[x] + y
	c.v[0xF] = 0
	if vx < c.v[x] {
		c.v[0xF] = 1
	}
	c.v[x] = vx
	c.pc += 2
}

// 8xy5 - SUB Vx, Vy
// Set Vx = Vx - Vy, set VF = NOT borrow.
//
// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
func (c *chip8) subVxVy(x, y uint8) {
	c.print("SUB V%01X, V%01X", x, y)
	c.v[0xF] = 0
	if c.v[x] > c.v[y] {
		c.v[0xF] = 1
	}
	c.v[x] -= c.v[y]
	c.pc += 2
}

// 8xy6 - SHR Vx {, Vy}
// Set Vx = Vx SHR 1.
//
// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
func (c *chip8) shrVx(x uint8) {
	c.print("SHR V%01X", x)
	c.v[0xF] = 0
	if c.v[x]&1 == 1 {
		c.v[0xF] = 1
	}
	c.v[x] = c.v[x] >> 1
}

// 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.
//
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
func (c *chip8) subnVxVy(x, y uint8) {
	c.print("SUBN V%01X, V%01X", x, y)
	c.v[0xF] = 0
	if c.v[y] > c.v[x] {
		c.v[0xF] = 1
	}
	c.v[x] = c.v[y] - c.v[x]
	c.pc += 2
}

// 8xyE - SHL Vx {, Vy}
// Set Vx = Vx SHL 1.
//
// If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
func (c *chip8) shlVx(x uint8) {
	c.print("SHL V%01X", x)
	c.v[0xF] = 0
	if c.v[x]&0x80 == 1 {
		c.v[0xF] = 1
	}
	c.v[x] = c.v[x] << 1
}

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
//
// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
func (c *chip8) sneVxVy(x, y uint8) {
	c.print("SNE V%01X, V%01X", x, y)
	if c.v[x] != c.v[y] {
		c.pc += 2
	}
	c.pc += 2
}

// Annn - LD I, addr
// Set I = nnn.
//
// The value of register I is set to nnn.
func (c *chip8) ldIAddr(addr uint16) {
	c.print("LD I, %04X", addr)
	c.i = addr
	c.pc += 2
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
//
// The program counter is set to nnn plus the value of V0.
func (c *chip8) jpV0Addr(addr uint16) {
	c.print("JP V0, %04X", addr)
	c.pc = addr + uint16(c.v[0])
}

// Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
//
// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.
func (c *chip8) rnd(x, b uint8) {
	c.print("RND V%01X, %d", x, b)
	c.v[x] = uint8(rand.Uint32()%256) & b
}

// Dxyn - DRW Vx, Vy, nibble
// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
//
// The interpreter reads n bytes from memory, starting at the address stored in I.
// These bytes are then displayed as sprites on screen at coordinates (Vx, Vy).
// Sprites are XORed onto the existing screen.
// If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0.
// If the sprite is positioned so part of it is outside the coordinates of the display,
// it wraps around to the opposite side of the screen.
// See instruction 8xy3 for more information on XOR, and section 2.4, Display, for more
// information on the Chip-8 screen and sprites.
func (c *chip8) drwVxVyN(x, y, n uint16) {
	c.print("DRW V%01X, %01X, %d", x, y, n)
	// TODO
	c.v[0xF] = 0
	c.pc += 2
}

// Ex9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
func (c *chip8) skpVx(x uint8) {
	c.print("SKP V%01X", x)
	if c.keypad[x] {
		c.pc += 2
	}
	c.pc += 2
}

// ExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is not pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
func (c *chip8) sknpVx(x uint8) {
	c.print("SKNP V%01X", x)
	if !c.keypad[x] {
		c.pc += 2
	}
	c.pc += 2
}

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
//
// The value of DT is placed into Vx.
func (c *chip8) ldVxDT(x uint8) {
	c.print("LD V%01X, DT", x)
	c.v[x] = c.dt
	c.pc += 2
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
//
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func (c *chip8) ldVxK(x uint8) {
	c.print("LD V%01X, K", x)
	// TODO
	var k uint8
	c.v[x] = k
	c.pc += 2
}

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
//
// DT is set equal to the value of Vx.
func (c *chip8) ldDTVx(x uint8) {
	c.print("LD DT, V%01X", x)
	c.dt = c.v[x]
	c.pc += 2
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
//
// ST is set equal to the value of Vx.
func (c *chip8) ldSTVx(x uint8) {
	c.print("LD ST, V%01X", x)
	c.st = c.v[x]
	c.pc += 2
}

// Fx1E - ADD I, Vx
// Set I = I + Vx.
//
// The values of I and Vx are added, and the results are stored in I.
func (c *chip8) addIVx(x uint8) {
	c.print("ADD I, V%01X", x)
	c.i += uint16(c.v[x])
	c.pc += 2
}

// Fx29 - LD F, Vx
// Set I = location of sprite for digit Vx.
//
// The value of I is set to the location for the hexadecimal sprite
// corresponding to the value of Vx. See section 2.4, Display, for more
// information on the Chip-8 hexadecimal font.
func (c *chip8) ldFVx(x uint8) {
	c.print("LD F, V%01X", x)
	// TODO
	c.pc += 2
}

// Fx33 - LD B, Vx
// Store BCD representation of Vx in memory locations I, I+1, and I+2.
//
// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
func (c *chip8) ldBVx(x uint8) {
	c.print("LD B, V%01X", x)
	c.pc += 2
}

// Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
//
// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
func (c *chip8) ldIVx(x uint8) {
	c.print("LD [I], V%01X", x)
	for i := uint8(0); i <= x; i++ {
		c.ram[c.i+uint16(i)] = c.v[x]
	}
	c.pc += 2
}

// Fx65 - LD Vx, [I]
// Read registers V0 through Vx from memory starting at location I.
//
// The interpreter reads values from memory starting at location I into registers V0 through Vx.
func (c *chip8) ldVxI(x uint8) {
	c.print("LD V%01X, [I]", x)
	for i := uint8(0); i <= x; i++ {
		c.v[x] = c.ram[c.i+uint16(i)]
	}
	c.pc += 2
}
