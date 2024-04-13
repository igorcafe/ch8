package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	var step bool
	var refreshPeriod time.Duration
	flag.DurationVar(&refreshPeriod, "r", 200*time.Microsecond, "refresh period duration")
	flag.BoolVar(&step, "step", false, "")
	flag.Parse()

	c8 := newChip8()
	log.SetFlags(0)

	b, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	copy(c8.ram[0x200:], b)

	scr, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = scr.Init()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			// cancel()
			scr.Fini()
			panic(r)
		}
	}()

	scr.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack))

	setText := func(x, y int, txt string, style tcell.Style) {
		for i, r := range txt {
			scr.SetContent(x+i, y, r, nil, style)
		}
	}

	events := make(chan tcell.Event, 0)
	quit := make(chan struct{}, 0)

	go scr.ChannelEvents(events, quit)

	c8.waitKey = func() {
		for {
			ev := <-events
			if _, ok := ev.(*tcell.EventKey); ok {
				break
			}
			events <- ev
		}
	}

	go func() {
		for {
			ev := <-events
			switch ev := ev.(type) {
			case *tcell.EventResize:
				scr.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					scr.Fini()
					return
				}

				for i := range c8.keypad {
					c8.keypad[i] = true
				}
				// if ev.Rune() == ' ' {
				// 	step = !step
				// }
			}
		}
	}()

	for {
		now := time.Now()

		c8.step()
		for step {
			ev := <-events
			if _, ok := ev.(*tcell.EventKey); ok {
				break
			}
			events <- ev
		}

		for lin := range c8.screen {
			for col := range c8.screen[lin] {
				if c8.screen[lin][col] {
					scr.SetContent(col, lin, ' ', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
				} else {
					scr.SetContent(col, lin, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
				}
			}
		}

		in := parseOpcode(c8.fetch(c8.pc))
		if in.id != "" {
			setText(83, 2, strings.Repeat(" ", 20), tcell.StyleDefault.Foreground(tcell.ColorGreenYellow))
			setText(83, 2, in.asm, tcell.StyleDefault.Foreground(tcell.ColorGreenYellow))
		}

		for x := 0; x <= 0xf; x++ {
			setText(90, 4+2*x, fmt.Sprintf("V%1X: %02X", x, c8.v[x]), tcell.StyleDefault)
		}
		for x := 0; x <= 0xf; x++ {
			k := 0
			if c8.keypad[x] {
				k = 1
			}
			setText(83, 4+2*x, fmt.Sprintf("K%1X: %1X", x, k), tcell.StyleDefault)
		}
		setText(98, 4, fmt.Sprintf("PC: %04X", c8.pc), tcell.StyleDefault)
		setText(98, 6, fmt.Sprintf("I:   %03X", c8.i), tcell.StyleDefault)
		setText(98, 8, fmt.Sprintf("RET: %03X", c8.stack[c8.sp]), tcell.StyleDefault)
		setText(98, 10, fmt.Sprintf("DT:  %02X", c8.dt), tcell.StyleDefault)
		setText(98, 12, fmt.Sprintf("ST:  %02X", c8.st), tcell.StyleDefault)
		setText(98, 14, fmt.Sprintf("[I]: %03X", c8.ram[c8.i]), tcell.StyleDefault)
		setText(98, 16, fmt.Sprintf("[PC]: %04X", c8.fetch(c8.pc)), tcell.StyleDefault)

		scr.Show()
		// c.drawToTerminal()
		time.Sleep(refreshPeriod - time.Since(now))
	}
}

// if you blur your vision, you'll see it a little better
var fontSet = []byte{
	// 0
	0b11110000,
	0b10010000,
	0b10010000,
	0b10010000,
	0b11110000,

	// 1
	0b00100000,
	0b01100000,
	0b00100000,
	0b00100000,
	0b01110000,

	// 2
	0b11110000,
	0b00010000,
	0b11110000,
	0b10000000,
	0b11110000,

	// 3
	0b11110000,
	0b00010000,
	0b11110000,
	0b00010000,
	0b11110000,

	// 4
	0b10010000,
	0b10010000,
	0b11110000,
	0b00010000,
	0b00010000,

	// 5
	0b11110000,
	0b10000000,
	0b11110000,
	0b00010000,
	0b11110000,

	// 6
	0b11110000,
	0b10000000,
	0b11110000,
	0b10010000,
	0b11110000,

	// 7
	0b11110000,
	0b00010000,
	0b00100000,
	0b01000000,
	0b01000000,

	// 8
	0b11110000,
	0b10010000,
	0b11110000,
	0b10010000,
	0b11110000,

	// 9
	0b11110000,
	0b10010000,
	0b11110000,
	0b00010000,
	0b11110000,

	// A
	0b11110000,
	0b10010000,
	0b11110000,
	0b10010000,
	0b10010000,

	// B
	0b11100000,
	0b10010000,
	0b11100000,
	0b10010000,
	0b11100000,

	// C
	0b11110000,
	0b10000000,
	0b10000000,
	0b10000000,
	0b11110000,

	// D
	0b11100000,
	0b10010000,
	0b10010000,
	0b10010000,
	0b11100000,

	// E
	0b11110000,
	0b10000000,
	0b11110000,
	0b10000000,
	0b11110000,

	// F
	0b11110000,
	0b10000000,
	0b11110000,
	0b10000000,
	0b10000000,
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
	screen [32][64]bool

	isKeyDown func(k uint8) bool
	waitKey   func()
}

func (c *chip8) fetch(pc uint16) uint16 {
	hi, lo := c.ram[pc], c.ram[pc+1]
	return uint16(hi)<<8 | uint16(lo)
}

func (c *chip8) step() {
	op := c.fetch(c.pc)
	in := parseOpcode(op)
	c.pc += 2
	if in.id == "" {
		log.Panicf("unknown opcode %04x, skipping\n", op)
		c.pc += 2
		return
	}

	// log.Print()
	c.print("\033[1;33m" + in.asm + "\033[0m")

	switch in.id {
	default:
		log.Panicf("unknown instruction: %04X, %#v", op, in)
	case "CLS":
		c.cls()
	case "RET":
		c.ret()
	case "SYS addr":
		c.sysAddr(in.addr)
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
	case "LD Vx, Vy":
		c.ldVxVy(in.x, in.y)
	case "OR Vx, Vy":
		c.orVxVy(in.x, in.y)
	case "AND Vx, Vy":
		c.andVxVy(in.x, in.y)
	case "XOR Vx, Vy":
		c.xorVxVy(in.x, in.y)
	case "ADD Vx, Vy":
		c.addVxVy(in.x, in.y)
	case "SUB Vx, Vy":
		c.subVxVy(in.x, in.y)
	case "SHR Vx {, Vy}":
		c.shrVx(in.x)
	case "SUBN Vx, Vy":
		c.subnVxVy(in.x, in.y)
	case "SHL Vx {, Vy}":
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
	if c.dt > 0 {
		c.dt--
	}
}

func (c *chip8) drawToTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()

	for x := range c.screen {
		for _, on := range c.screen[x] {
			if on {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
