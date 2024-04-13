package main

import (
	"fmt"
	"log"
	"math/rand"
)

type instruction struct {
	id   string
	asm  string
	x    uint8
	y    uint8
	b    uint8
	n    uint8
	addr uint16
}

func parseOpcode(op uint16) instruction {
	// 00E0 - CLS
	if op == 0x00E0 {
		return instruction{
			id:  "CLS",
			asm: "CLS",
		}
	}

	// 00EE - RET
	if op == 0x00EE {
		return instruction{
			id:  "RET",
			asm: "RET",
		}
	}

	// 0nnn - SYS addr
	if op&0xF000 == 0x0000 {
		addr := op & 0x0FFF
		return instruction{
			id:   "SYS addr",
			asm:  fmt.Sprintf("SYS %04X", addr),
			addr: addr,
		}
	}

	// 1nnn - JP addr
	if op&0xF000 == 0x1000 {
		addr := op & 0x0FFF
		return instruction{
			id:   "JP addr",
			asm:  fmt.Sprintf("JP %04X", addr),
			addr: addr,
		}
	}

	// 2nnn - CALL addr
	if op&0xF000 == 0x2000 {
		addr := op & 0x0FFF
		return instruction{
			id:   "CALL addr",
			asm:  fmt.Sprintf("CALL %04X", addr),
			addr: addr,
		}
	}

	// 3xkk - SE Vx, byte
	if op&0xF000 == 0x3000 {
		x := uint8((op & 0x0F00) >> 8)
		b := uint8(op & 0x00FF)
		return instruction{
			id:  "SE Vx, byte",
			asm: fmt.Sprintf("SE V%01X, %02X", x, b),
			x:   x,
			b:   b,
		}
	}

	// 4xkk - SNE Vx, byte
	if op&0xF000 == 0x4000 {
		x := uint8((op >> 8) & 0xF)
		b := uint8(op & 0xFF)
		return instruction{
			id:  "SNE Vx, byte",
			asm: fmt.Sprintf("SNE V%01X, %02X", x, b),
			x:   x,
			b:   b,
		}
	}

	// 5xy0 - SE Vx, Vy
	if op&0xF000 == 0x5000 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "SE Vx, Vy",
			asm: fmt.Sprintf("SE V%01X, V%01X", x, y),
		}
	}

	// 6xkk - LD Vx, byte
	if op&0xF000 == 0x6000 {
		x := uint8((op >> 8) & 0xF)
		b := uint8(op & 0xFF)
		return instruction{
			id:  "LD Vx, byte",
			asm: fmt.Sprintf("LD V%01X, %02X", x, b),
			x:   x,
			b:   b,
		}
	}

	// 7xkk - ADD Vx, byte
	if op&0xF000 == 0x7000 {
		x := uint8((op >> 8) & 0xF)
		b := uint8(op & 0xFF)
		return instruction{
			id:  "ADD Vx, byte",
			asm: fmt.Sprintf("ADD V%01X, %02X", x, b),
			x:   x,
			b:   b,
		}
	}

	// 8xy0 - LD Vx, Vy
	if op&0xF00F == 0x8000 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "LD Vx, Vy",
			asm: fmt.Sprintf("LD V%01X, V%01X", x, y),
			x:   x,
			y:   y,
		}
	}

	// 8xy1 - OR Vx, Vy
	if op&0xF00F == 0x8001 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "OR Vx, Vy",
			asm: fmt.Sprintf("OR V%01X, V%01X", x, y),
			x:   x,
			y:   y,
		}
	}

	// 8xy2 - AND Vx, Vy
	if op&0xF00F == 0x8002 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "AND Vx, Vy",
			asm: fmt.Sprintf("AND V%01X, V%01X", x, y),
			x:   x,
			y:   y,
		}
	}

	// 8xy3 - XOR Vx, Vy
	if op&0xF00F == 0x8003 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "XOR Vx, Vy",
			asm: fmt.Sprintf("XOR V%01X, V%01X", x, y),
			x:   x,
			y:   y,
		}
	}

	// 8xy4 - ADD Vx, Vy
	if op&0xF00F == 0x8004 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "ADD Vx, Vy",
			asm: fmt.Sprintf("ADD V%01X, V%01X", x, y),
			x:   x,
			y:   y,
		}
	}

	// 8xy5 - SUB Vx, Vy
	if op&0xF00F == 0x8005 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "SUB Vx, Vy",
			asm: fmt.Sprintf("SUB V%01X, V%01X", x, y),
			x:   x,
			y:   y,
		}
	}

	// 8xy6 - SHR Vx
	if op&0xF00F == 0x8006 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "SHR Vx {, Vy}",
			asm: fmt.Sprintf("SUB V%01X, {, V%01X}", x, y),
			x:   x,
			y:   y,
		}
	}

	// 8xy7 - SUBN Vx, Vy
	if op&0xF00F == 0x8007 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "SUBN Vx, Vy",
			asm: fmt.Sprintf("SUBN V%01X, V%01X", x, y),
			x:   x,
			y:   y,
		}
	}

	// 8xyE - SHL Vx{, Vy}
	if op&0xF00F == 0x800E {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "SHL Vx {, Vy}",
			asm: fmt.Sprintf("SHL V%01X {, V%01X}", x, y),
			x:   x,
			y:   y,
		}
	}

	// 9xy0 - SNE Vx, Vy
	if op&0xF000 == 0x9000 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		return instruction{
			id:  "SNE Vx, Vy",
			asm: fmt.Sprintf("SNE V%01X, V%01X", x, y),
			x:   x,
			y:   y,
		}
	}

	// Annn - LD I, addr
	if op&0xF000 == 0xA000 {
		addr := op & 0xFFF
		return instruction{
			id:   "LD I, addr",
			asm:  fmt.Sprintf("LD I, %04X", addr),
			addr: addr,
		}
	}

	// Bnnn - JP V0, addr
	if op&0xF000 == 0xA000 {
		addr := op & 0xFFF
		return instruction{
			id:   "JP V0, addr",
			asm:  fmt.Sprintf("JP V0, %04X", addr),
			addr: addr,
		}
	}

	// Cxkk - RND Vx, byte
	if op&0xF000 == 0xC000 {
		x := uint8((op >> 8) & 0xF)
		b := uint8(op & 0xFF)
		return instruction{
			id:  "RND Vx, byte",
			asm: fmt.Sprintf("RND V%01X, %02X", x, b),
			x:   x,
			b:   b,
		}
	}

	// Dxyn - DRW, Vx, Vy, nibble
	if op&0xF000 == 0xD000 {
		x := (op & 0x0F00) >> 8
		y := (op & 0x00F0) >> 4
		n := op & 0x000F
		return instruction{
			id:  "DRW Vx, Vy, nibble",
			asm: fmt.Sprintf("DRW V%01X, V%01X, %d", x, y, n),
			x:   uint8(x),
			y:   uint8(y),
			n:   uint8(n),
		}
	}

	// Ex9E - SKP Vx
	if op&0xF00F == 0xE00E {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "SKP Vx",
			asm: fmt.Sprintf("SKP V%01X", x),
			x:   x,
		}
	}

	// ExA1 - SKNP Vx
	if op&0xF00F == 0xE001 {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "SKNP Vx",
			asm: fmt.Sprintf("SKNP V%01X", x),
			x:   x,
		}
	}

	// Fx07 - LD Vx, DT
	if op&0xF0FF == 0xF007 {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "LD Vx, DT",
			asm: fmt.Sprintf("LD V%01X, DT", x),
			x:   x,
		}
	}

	// Fx0A - LD Vx, K
	if op&0xF0FF == 0xF00A {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "LD Vx, K",
			asm: fmt.Sprintf("LD V%01X, K", x),
			x:   x,
		}
	}

	// Fx15 - LD DT, Vx
	if op&0xF0FF == 0xF015 {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "LD DT, Vx",
			asm: fmt.Sprintf("LD DT, V%01X", x),
			x:   x,
		}
	}

	// Fx18 - LD ST, Vx
	if op&0xF0FF == 0xF018 {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "LD ST, Vx",
			asm: fmt.Sprintf("LD ST, V%01X", x),
			x:   x,
		}
	}

	// Fx1E - ADD I, Vx
	if op&0xF0FF == 0xF01E {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "ADD I, Vx",
			asm: fmt.Sprintf("ADD I, V%01X", x),
			x:   x,
		}
	}

	// Fx29 - LD F, Vx
	if op&0xF0FF == 0xF029 {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "LD F, Vx",
			asm: fmt.Sprintf("LD F, V%01X", x),
			x:   x,
		}
	}

	// Fx33 - LD B, Vx
	if op&0xF0FF == 0xF033 {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "LD B, Vx",
			asm: fmt.Sprintf("LD B, V%01X", x),
			x:   x,
		}
	}

	// Fx55 - LD [I], Vx
	if op&0xF0FF == 0xF055 {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "LD [I], Vx",
			asm: fmt.Sprintf("LD [I], V%01X", x),
			x:   x,
		}
	}

	// Fx65 - LD Vx, [I]
	if op&0xF0FF == 0xF065 {
		x := uint8((op >> 8) & 0xF)
		return instruction{
			id:  "LD Vx, [I]",
			asm: fmt.Sprintf("LD V%01X, [I]", x),
			x:   x,
		}
	}

	return instruction{}
}

func (c *chip8) print(format string, args ...any) {
	return
	log.Printf("\033[1;33m"+format+"\033[0m", args...)
	log.Printf("| PC: %04X |  I: %04X | RET: %03X | DT:   %02X | ST:   %02X | [I]: %02X", c.pc, c.i, c.stack[c.sp], c.dt, c.st, c.ram[c.i])
	s := "|"
	for i := 0; i < 8; i++ {
		s += fmt.Sprintf(" V%01X:   %02X |", i, c.v[i])
	}
	s += "\n|"
	for i := 8; i <= 0xF; i++ {
		s += fmt.Sprintf(" V%01X:   %02X |", i, c.v[i])
	}
	log.Print(s)
}

// 0nnn - SYS addr
// Jump to a machine code routine at nnn.
//
// This instruction is only used on the old computers on which Chip-8
// was originally implemented. It is ignored by modern interpreters.
func (c *chip8) sysAddr(addr uint16) {
	c.pc = addr
}

// 00E0 - CLS
// Clear the display.
func (c *chip8) cls() {
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
	c.pc = c.stack[c.sp]
	c.sp--
	c.pc += 2
}

// 1nnn - JP addr
// Jump to location nnn.
//
// The interpreter sets the program counter to nnn.
func (c *chip8) jpAddr(addr uint16) {
	c.pc = addr
}

// 2nnn - CALL addr
// Call subroutine at nnn.
//
// The interpreter increments the stack pointer, then puts the current PC
// on the top of the stack. The PC is then set to nnn.
func (c *chip8) callAddr(addr uint16) {
	c.sp++
	c.stack[c.sp] = c.pc
	c.pc = addr
}

// 3xkk - SE Vx, byte
// Skip next instruction if Vx = kk.
//
// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
func (c *chip8) seVxB(x, b uint8) {
	if c.v[x] == b {
		c.pc += 2
	}
	c.pc += 2
}

// 4xkk - SNE Vx, byte
// Skip next instruction if Vx != kk.
//
// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
func (c *chip8) sneVxB(x, b uint8) {
	if c.v[x] != b {
		c.pc += 2
	}
	c.pc += 2
}

// 5xy0 - SE Vx, Vy
// Skip next instruction if Vx = Vy.
//
// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
func (c *chip8) seVxVy(x, y uint8) {
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
	c.v[x] += b
	c.pc += 2
}

// 8xy0 - LD Vx, Vy
// Set Vx = Vy.
//
// Stores the value of register Vy in register Vx.
func (c *chip8) ldVxVy(x, y uint8) {
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
	c.v[x] |= c.v[y]
	c.pc += 2
}

// 8xy2 - AND Vx, Vy
// Set Vx = Vx AND Vy.
//
// Performs a bitwise AND on the values of Vx andVxVy Vy, then stores the result in Vx. A bitwise AND compares the corrseponding bits from two values, andVxVy if both bits are 1, then the same bit in the result is also 1. Otherwise, it is 0.
func (c *chip8) andVxVy(x, y uint8) {
	c.v[x] &= c.v[y]
	c.pc += 2
}

// 8xy3 - XOR Vx, Vy
// Set Vx = Vx XOR Vy.
//
// Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx. An exclusive OR compares the corrseponding bits from two values, and if the bits are not both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
func (c *chip8) xorVxVy(x, y uint8) {
	c.v[x] ^= c.v[y]
	c.pc += 2
}

// 8xy4 - ADD Vx, Vy
// Set Vx = Vx + Vy, set VF = carry.
//
// The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
func (c *chip8) addVxVy(x, y uint8) {
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
	c.v[0xF] = 0
	if c.v[x]&1 == 1 {
		c.v[0xF] = 1
	}
	c.v[x] = c.v[x] >> 1
	c.pc += 2
}

// 8xy7 - SUBN Vx, Vy
// Set Vx = Vy - Vx, set VF = NOT borrow.
//
// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
func (c *chip8) subnVxVy(x, y uint8) {
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
	c.v[0xF] = 0
	if c.v[x]&0x80 == 1 {
		c.v[0xF] = 1
	}
	c.v[x] = c.v[x] << 1
	c.pc += 2
}

// 9xy0 - SNE Vx, Vy
// Skip next instruction if Vx != Vy.
//
// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
func (c *chip8) sneVxVy(x, y uint8) {
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
	c.i = addr
	c.pc += 2
}

// Bnnn - JP V0, addr
// Jump to location nnn + V0.
//
// The program counter is set to nnn plus the value of V0.
func (c *chip8) jpV0Addr(addr uint16) {
	c.pc = addr + uint16(c.v[0])
}

// Cxkk - RND Vx, byte
// Set Vx = random byte AND kk.
//
// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.
func (c *chip8) rndVxB(x, b uint8) {
	c.v[x] = uint8(rand.Uint32()%256) & b
	c.pc += 2
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
func (c *chip8) drwVxVyN(x, y, n uint8) {
	c.v[0xF] = 0
	for i := uint8(0); i < n; i++ {
		lin := c.v[y] + i
		b := c.ram[c.i+uint16(i)]
		for col := c.v[x]; col < c.v[x]+8; col++ {
			set := b>>(7-col%8)&1 == 1
			if c.screen[lin][col] && set {
				c.v[0xF] = 1
				set = false
			}
			c.screen[lin][col] = set
		}
	}
	c.pc += 2
}

// Ex9E - SKP Vx
// Skip next instruction if key with the value of Vx is pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
func (c *chip8) skpVx(x uint8) {
	if c.keypad[c.v[x]] {
		c.pc += 2
	}
	c.pc += 2
}

// ExA1 - SKNP Vx
// Skip next instruction if key with the value of Vx is not pressed.
//
// Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
func (c *chip8) sknpVx(x uint8) {
	if !c.keypad[c.v[x]] {
		c.pc += 2
	}
	c.pc += 2
}

// Fx07 - LD Vx, DT
// Set Vx = delay timer value.
//
// The value of DT is placed into Vx.
func (c *chip8) ldVxDT(x uint8) {
	c.v[x] = c.dt
	c.pc += 2
}

// Fx0A - LD Vx, K
// Wait for a key press, store the value of the key in Vx.
//
// All execution stops until a key is pressed, then the value of that key is stored in Vx.
func (c *chip8) ldVxK(x uint8) {
	panic("todo")
	var k uint8
	c.v[x] = k
	c.pc += 2
}

// Fx15 - LD DT, Vx
// Set delay timer = Vx.
//
// DT is set equal to the value of Vx.
func (c *chip8) ldDTVx(x uint8) {
	c.dt = c.v[x]
	c.pc += 2
}

// Fx18 - LD ST, Vx
// Set sound timer = Vx.
//
// ST is set equal to the value of Vx.
func (c *chip8) ldSTVx(x uint8) {
	c.st = c.v[x]
	c.pc += 2
}

// Fx1E - ADD I, Vx
// Set I = I + Vx.
//
// The values of I and Vx are added, and the results are stored in I.
func (c *chip8) addIVx(x uint8) {
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
	panic("todo")
	c.pc += 2
}

// Fx33 - LD B, Vx
// Store BCD representation of Vx in memory locations I, I+1, and I+2.
//
// The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
func (c *chip8) ldBVx(x uint8) {
	c.pc += 2
}

// Fx55 - LD [I], Vx
// Store registers V0 through Vx in memory starting at location I.
//
// The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
func (c *chip8) ldIVx(x uint8) {
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
	for i := uint8(0); i <= x; i++ {
		c.v[x] = c.ram[c.i+uint16(i)]
	}
	c.pc += 2
}
