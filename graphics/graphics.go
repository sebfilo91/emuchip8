// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import "github.com/nsf/termbox-go"

//const chars = "*"

var output_mode = termbox.OutputNormal

func Init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
}

func Draw(sx int, sy int, char rune) {
	termbox.SetCell(sx, sy, char, termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
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