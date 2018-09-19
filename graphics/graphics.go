// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"github.com/nsf/termbox-go"
	"fmt"
)

//const chars = "*"

var output_mode = termbox.OutputNormal

const screenWidth = 64
const screenHeight = 32

var pixels [screenWidth * screenHeight]byte

var debugEnabled = false

func Init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
}

func Draw(x int, y int, sprite []byte) bool {
	n := len(sprite)
	collision := false

	for yline := 0; yline < n; yline++ {
		pixelRow := sprite[yline]

		for xline := 0; xline < 8; xline++ {

			// Pixel is set to one
			if(pixelRow & (0x80 >> byte(xline)) != 0) {
				if(pixels[(x + xline + ((y + yline) * 64))] == 1) {
	          		collision = true;                                 
				}

	        	pixels[x + xline + ((y + yline) * 64)] = pixels[x + xline + ((y + yline) * 64)] ^ 1;
			} 
		}
	}
	debug("\n --DRAW --")
	for i := 0; i < len(pixels); i++ {
		if(pixels[i] != 0) {
			debug("%v ", pixels[i])
		}
	}

	return collision
}

func Render() {
	debug("\n --RENDER --")
	for i := 0; i < len(pixels); i++ {
		if(pixels[i] != 0) {
			debug("%v ", pixels[i])
		}
	}
	for x := 0; x < screenWidth; x++ {
		for y := 0; y < screenHeight; y++ {
			if(pixels[(y*screenWidth) + x] == 0x1) {
				termbox.SetCell(x, y, rune('*'), termbox.ColorDefault, termbox.ColorDefault)
				termbox.Flush()
			}
		}
	}
}


func debug(msg string, object ...interface{}) {
	if debugEnabled == true {
		fmt.Printf(msg, object...)
	}
}

/*func draw_all() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	switch output_mode {

	case termbox.OutputNormal:
		print_one(12, 12)
	}

	termbox.Flush()
}*/


/*func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	draw_all()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
				case termbox.KeyEsc:
					break loop
			}
		case termbox.EventResize:
			draw_all()
		}
	}
}*/