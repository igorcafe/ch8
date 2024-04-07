package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
)

func main() {
	var step bool
	flag.BoolVar(&step, "step", false, "")
	flag.Parse()

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
		if step {
			scanner.Scan()
		}

		// c.drawToTerminal()
		// time.Sleep(time.Millisecond)
	}

	scanner.Scan()
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

	in := parseOpcode(op)
	if in.id == "" {
		fmt.Printf("unknown opcode %04x, skipping\n", op)
		c.pc += 2
		return
	}

	log.Print("\033[1;33m" + in.asm + "\033[0m")

	switch in.id {
	default:
		log.Panicf("unknown instruction: %04X, %#v", op, in)
	case "SYS addr":
		c.sysAddr(in.addr)
	case "CLS":
		c.cls()
	case "RET":
		c.ret()
	case "JP addr":
		c.jpAddr(in.addr)
	case "CALL addr":
		c.callAddr(in.addr)
	case "SE Vx, byte":
		c.seVxB(in.x, in.b)
	case "SNE Vx, byte":
		c.sneVxB(in.x, in.b)
	case "SE Vx, Vy":
		c.seVxVy(in.x, in.y)
	case "LD Vx, byte":
		c.ldVxB(in.x, in.b)
	case "ADD Vx, byte":
		c.addVxB(in.x, in.b)
	case "OR Vx, Vy":
		c.orVxVy(in.x, in.y)
	case "AND Vx, Vy":
		c.andVxVy(in.x, in.y)
	case "ADD Vx, Vy":
		c.addVxVy(in.x, in.y)
	case "SUB Vx, Vy":
		c.subVxVy(in.x, in.y)
	case "SHR Vx":
		c.shrVx(in.x)
	case "SUBN Vx, Vy":
		c.subnVxVy(in.x, in.y)
	case "SHL Vx":
		c.shlVx(in.x)
	case "SNE Vx, Vy":
		c.sneVxVy(in.x, in.y)
	case "LD I, addr":
		c.ldIAddr(in.addr)
	case "JP V0, addr":
		c.jpV0Addr(in.addr)
	case "RND Vx, byte":
		c.rndVxB(in.x, in.b)
	case "DRW Vx, Vy, nibble":
		c.drwVxVyN(in.x, in.y, in.n)
	case "SKP Vx":
		c.skpVx(in.x)
	case "SKNP Vx":
		c.sknpVx(in.x)
	case "LD Vx, DT":
		c.ldVxDT(in.x)
	case "LD Vx, K":
		c.ldVxK(in.x)
	case "LD DT, Vx":
		c.ldDTVx(in.x)
	case "LD ST, Vx":
		c.ldSTVx(in.x)
	case "ADD I, Vx":
		c.addIVx(in.x)
	case "LD F, Vx":
		c.ldFVx(in.x)
	case "LD B, Vx":
		c.ldBVx(in.x)
	case "LD [I], Vx":
		c.ldIVx(in.x)
	case "LD Vx, [I]":
		c.ldVxI(in.x)
	}
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
