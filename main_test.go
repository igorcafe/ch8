package main

import (
	"math/rand"
	"testing"
)

func Test_chip8_cls(t *testing.T) {
	c8 := newChip8()
	empty := c8.screen

	// fill screen with random data
	for x := range c8.screen {
		for y := range c8.screen[x] {
			c8.screen[x][y] = rand.Float64() > 0.5
		}
	}

	c8.cls()
	if c8.screen != empty {
		t.Fatalf("want: empty screen, got: %#v", c8.screen)
	}
}
