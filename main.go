package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
)

func main() {
	c := newChip8()
	log.SetFlags(0)

	for i := range c.screen {
		for j := range c.screen[i] {
			c.screen[i][j] = rand.Float64() > 0.5
		}
	}

	b, err := os.ReadFile("invaders.ch8")
	if err != nil {
		panic(err)
	}

	copy(c.ram[0x200:], b)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		c.step()
		scanner.Scan()

		// c.drawToTerminal()
		// time.Sleep(time.Millisecond)
	}
}

var fontSet = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func newChip8() *chip8 {
	c := &chip8{
		pc: 0x200,
	}
	copy(c.ram[0:80], fontSet)
	return c
}

type chip8 struct {
	// general purpose registers
	v [16]uint8

	// I register (?)
	i uint16

	// delay timer (DT)
	dt uint8

	// sound timer (ST)
	st uint8

	// program counter
	pc uint16

	// stack pointer
	sp uint16

	// stack where sp points to
	stack [16]uint16

	// keypad represents the state of the keys
	keypad [16]bool

	// emulated ram
	ram [4096]uint8

	// state of screen per pixel (on/off)
	screen [64][32]bool
}

func (c *chip8) step() {
	hi, lo := c.ram[c.pc], c.ram[c.pc+1]
	op := uint16(hi)<<8 | uint16(lo)

	// log.Printf("op: %04X", op)

	// 00E0 - CLS
	if op == 0x00E0 {
		c.cls()
		return
	}

	// 00EE - RET
	if op == 0x00EE {
		c.ret()
		return
	}

	// 0nnn - SYS addr
	if op&0xF000 == 0x0000 {
		addr := op & 0x0FFF
		c.sysAddr(addr)
		return
	}

	// 1nnn - JP addr
	if op&0xF000 == 0x1000 {
		addr := op & 0x0FFF
		c.jpAddr(addr)
		return
	}

	// 2nnn - CALL addr
	if op&0xF000 == 0x2000 {
		addr := op & 0x0FFF
		c.callAddr(addr)
		return
	}

	// 3xkk - SE Vx, byte
	if op&0xF000 == 0x3000 {
		x := uint8((op & 0x0F00) >> 8)
		b := uint8(op & 0x00FF)
		c.seVxB(x, b)
		return
	}

	// 4xkk - SNE Vx, byte
	if op&0xF000 == 0x4000 {
		x := uint8((op >> 8) & 0xF)
		b := uint8(op & 0xFF)
		c.sneVxB(x, b)
		return
	}

	// 5xy0 - SE Vx, Vy
	if op&0xF000 == 0x5000 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		c.seVxVy(x, y)
		return
	}

	// 6xkk - LD Vx, byte
	if op&0xF000 == 0x6000 {
		x := uint8((op >> 8) & 0xF)
		b := uint8(op & 0xFF)
		c.ldVxB(x, b)
		return
	}

	// 7xkk - ADD Vx, byte
	if op&0xF000 == 0x7000 {
		x := uint8((op >> 8) & 0xF)
		b := uint8(op & 0xFF)
		c.addVxB(x, b)
		return
	}

	// 8xy1 - OR Vx, Vy
	if op&0xF00F == 0x8001 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		c.orVxVy(x, y)
		return
	}

	// 8xy2 - AND Vx, Vy
	if op&0xF00F == 0x8002 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		c.andVxVy(x, y)
		return
	}

	// 8xy3 - XOR Vx, Vy
	if op&0xF00F == 0x8003 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		c.xorVxVy(x, y)
		return
	}

	// 8xy4 - ADD Vx, Vy
	if op&0xF00F == 0x8004 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		c.addVxVy(x, y)
		return
	}

	// 8xy5 - SUB Vx, Vy
	if op&0xF00F == 0x8005 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		c.subVxVy(x, y)
		return
	}

	// 8xy6 - SHR Vx
	if op&0xF00F == 0x8006 {
		x := uint8((op >> 8) & 0xF)
		_ = uint8((op >> 4) & 0xF)
		c.shrVx(x)
		return
	}

	// 8xy7 - SUBN Vx, Vy
	if op&0xF00F == 0x8007 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		c.subnVxVy(x, y)
		return
	}

	// 8xyE - SHL Vx
	if op&0xF00F == 0x800E {
		x := uint8((op >> 8) & 0xF)
		_ = uint8((op >> 4) & 0xF)
		c.shlVx(x)
		return
	}

	// 9xy0 - SNE Vx, Vy
	if op&0xF000 == 0x9000 {
		x := uint8((op >> 8) & 0xF)
		y := uint8((op >> 4) & 0xF)
		c.sneVxVy(x, y)
		return
	}

	// Annn - LD I, addr
	if op&0xF000 == 0xA000 {
		addr := op & 0xFFF
		c.ldIAddr(addr)
		return
	}

	// Bnnn - JP V0, addr
	if op&0xF000 == 0xA000 {
		addr := op & 0xFFF
		c.jpV0Addr(addr)
		return
	}

	// Cxkk - RND Vx, byte
	if op&0xF000 == 0xA000 {
		x := uint8((op >> 8) & 0xF)
		b := uint8(op & 0xFF)
		c.rnd(x, b)
		return
	}

	// Dxyn - DRW, Vx, Vy, nibble
	if op&0xF000 == 0xD000 {
		x := (op & 0x0F00) >> 8
		y := (op & 0x00F0) >> 4
		n := op & 0x000F
		c.drwVxVyN(x, y, n)
		return
	}

	// Ex9E - SKP Vx
	if op&0xF00F == 0xE00E {
		x := uint8((op >> 8) & 0xF)
		c.skpVx(x)
		return
	}

	// ExA1 - SKNP Vx
	if op&0xF00F == 0xE001 {
		x := uint8((op >> 8) & 0xF)
		c.sknpVx(x)
		return
	}

	// Fx07 - LD Vx, DT
	if op&0xF0FF == 0xF007 {
		x := uint8((op >> 8) & 0xF)
		c.ldVxDT(x)
		return
	}

	// Fx0A - LD Vx, K
	if op&0xF0FF == 0xF00A {
		x := uint8((op >> 8) & 0xF)
		c.ldVxK(x)
		return
	}

	// Fx15 - LD DT, Vx
	if op&0xF0FF == 0xF015 {
		x := uint8((op >> 8) & 0xF)
		c.ldDTVx(x)
		return
	}

	// Fx18 - LD ST, Vx
	if op&0xF0FF == 0xF018 {
		x := uint8((op >> 8) & 0xF)
		c.ldSTVx(x)
		return
	}

	// Fx1E - ADD I, Vx
	if op&0xF0FF == 0xF01E {
		x := uint8((op >> 8) & 0xF)
		c.addIVx(x)
		return
	}

	// Fx29 - LD F, Vx
	if op&0xF0FF == 0xF029 {
		x := uint8((op >> 8) & 0xF)
		c.ldFVx(x)
		return
	}

	// Fx33 - LD B, Vx
	if op&0xF0FF == 0xF033 {
		x := uint8((op >> 8) & 0xF)
		c.ldBVx(x)
		return
	}

	// Fx55 - LD [I], Vx
	if op&0xF0FF == 0xF055 {
		x := uint8((op >> 8) & 0xF)
		c.ldIVx(x)
		return
	}

	// Fx65 - LD Vx, [I]
	if op&0xF0FF == 0xF065 {
		x := uint8((op >> 8) & 0xF)
		c.ldVxI(x)
		return
	}

	fmt.Printf("unknown opcode %04x, skipping\n", op)
	c.pc += 2
}

func (c *chip8) drawToTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()

	for x := range c.screen {
		for _, on := range c.screen[x] {
			if on {
				fmt.Print("##")
			} else {
				fmt.Print("..")
			}
		}
		fmt.Println()
	}
}
