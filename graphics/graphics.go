// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import "github.com/nsf/termbox-go"

//const chars = "*"

var output_mode = termbox.OutputNormal

const screenWidth = 64
const screenHeight = 32

var pixels [screenWidth * screenHeight]byte

func Init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
}

func Draw(sx int, sy int, sprite []byte) bool {
	n := len(sprite)
	collision := false

	for y := 0; y < n; y++ {
		pixelRow := sprite[y]

		for x := 0; x < 8; x++ {

			// Pixel is set to one
			if(pixelRow & (0x80 >> byte(x)) != 0) {

			} 
		}
	}


	//termbox.SetCell(sx, sy, char, termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()

	return collision
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