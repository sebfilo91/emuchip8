package main

import(
	"fmt"
	"math/rand"
    "io/ioutil"
    "time"
    "os"
    "github.com/nsf/termbox-go"
    graphics "emuchip8/graphics"
)

//Documentation http://mattmik.com/files/chip8/mastering/chip8.html

var debugEnabled = false

var opcode uint16

var memory [4096]byte

var stack[16]uint16

var SP uint16

var V [16]byte

var I uint16

// Program Counter
var PC uint16

// Delay Timer
var DT byte

// Sound Timer
var ST byte

var Fonts = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0,
	0x20, 0x60, 0x20, 0x20, 0x70,
	0xF0, 0x10, 0xF0, 0x80, 0xF0,
	0xF0, 0x10, 0xF0, 0x10, 0xF0,
	0x90, 0x90, 0xF0, 0x10, 0x10,
	0xF0, 0x80, 0xF0, 0x10, 0xF0,
	0xF0, 0x80, 0xF0, 0x90, 0xF0,
	0xF0, 0x10, 0x20, 0x40, 0x40,
	0xF0, 0x90, 0xF0, 0x90, 0xF0,
	0xF0, 0x90, 0xF0, 0x10, 0xF0,
	0xF0, 0x90, 0xF0, 0x90, 0x90,
	0xE0, 0x90, 0xE0, 0x90, 0xE0,
	0xF0, 0x80, 0x80, 0x80, 0xF0,
	0xE0, 0x90, 0x90, 0x90, 0xE0,
	0xF0, 0x80, 0xF0, 0x80, 0xF0,
	0xF0, 0x80, 0xF0, 0x80, 0x80,
}

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	loadFonts()
	loadRom("./roms/pong.rom")
	//loadRom("./roms/test.rom")

  	for {
		emulateCycle(0)
  	}
}

func loadFonts() {
	for i := 0; i < len(Fonts); i++ {
		memory[0x0 + i] = Fonts[i]
		//debug("\n0x%04X ", memory[loadAddress + i] )
	}
}

func loadRom(rom string) {
	// load at 0x200

	dat, err := ioutil.ReadFile(rom)
	if err != nil {
		debug("%e", err)
	}

	loadAddress := int(0x200)
	PC = uint16(loadAddress)

	for i := 0; i < len(dat); i++ {
		memory[loadAddress + i] = dat[i]
		//debug("\n0x%04X ", memory[loadAddress + i] )
	}

	debug("\n\nSTARTING MEMORY 0x%04X \n", memory[0x200])
	debug("\nSTARTING PC 0x%04X \n", PC)
	debug("\nSTARTING memory[PC] 0x%04X \n", memory[PC])
}

func initialize() {

}

func emulateCycle(ex int) {
	opcode = uint16(memory[PC]) << 8 | uint16(memory[PC + 1])
	debug("\nex (%d) | opcode : 0x%04X", ex, opcode)
	executeInstruction(opcode)
}

func executeInstruction(op uint16) {
	switch(op & 0xF000) {
		case 0x0000: 
			switch(op & 0x00FF) {
				case 0x00E0: 
					fmt.Println("Clear graphics")
					graphics.Clear()
					PC += 2
				case 0x00EE:
					SP--
					PC = stack[SP]
					PC += 2
				default:
					debug("Invalid opcode")
			}
		case 0x1000:
			PC = 0x0FFF & op
		case 0x2000:
			stack[SP] = PC
			SP++;
			PC = 0x0FFF & op
		case 0x3000:
			x := (op & 0x0F00) >> 8
			kk := (op & 0x00FF)

			PC += 2

			if (V[x] == byte(kk)) {
				PC += 2
			}
		case 0x4000:
			x := (op & 0x0F00) >> 8
			kk := (op & 0x00FF)

			PC += 2

			if (V[x] != byte(kk)) {
				PC += 2
			}
		case 0x5000:
			x := (op & 0x0F00) >> 8
			y := (op & 0x00F0) >> 4

			PC += 2

			if (V[x] ==  V[y]) {
				PC += 2
			}
		case 0x6000:
			x := (op & 0x0F00) >> 8
			kk := op & 0x00FF
			V[x] = byte(kk)

			PC += 2

		case 0x7000:
			addr := (op & 0x0F00) >> 8
			kk := byte(op)
			V[addr] = V[addr] + byte(kk)
			
			PC += 2
		case 0x8000:
			x := (op & 0x0F00) >> 8
			y := (op & 0x00F0) >> 4

			switch(op & 0xF00F) {
					
				case 0x8000:
					V[x] = V[y]

					PC += 2	
				case 0x8001:
					V[x] = V[x] | V[y]

					PC += 2	
				case 0x8002:
					V[x] = V[x] & V[y]

					PC += 2	
				case 0x8003:
					V[x] = V[x] ^ V[y]

					PC += 2	
				case 0x8004:
					result := uint16(V[x]) + uint16(V[y])

					carry := byte(0)

					if(result > 0xFF) {
						carry = 1
					}
					V[0xF] = carry
					V[x] = byte(result)

					PC += 2
				case 0x8005:
					result := uint16(V[x]) - uint16(V[y])

					V[0xF] = 0
					if V[x] > V[y] {
						V[0xF] = 1
					}

					V[x] = byte(result)

					PC += 2
				case 0x8006:
					/*
                    V[0xF] = V[(opcode & 0x0F00) >> 8] & 0x1;
                    V[(opcode & 0x0F00) >> 8] >>= 1;
					*/
					if (V[x] & 0x01) == 0x01 {
						V[0xF] = 1
					}

					V[x] = V[x] / 2

					PC += 2	

				case 0x8007:
					result := uint16(V[y]) - uint16(V[x])

					V[0xF] = 0
					if V[y] > V[x] {
						V[0xF] = 1
					}

					V[x] = byte(result)
					
					PC += 2

				case 0x800E:
					if (V[x] & 0x80) == 0x80 {
						V[0xF] = 1
					}

					V[x] = V[x] * 2

					PC += 2	
				default:
					debug("Invalid opcode")

			}

		case 0x9000:
			x := (op & 0x0F00) >> 8
			y := (op & 0x00F0) >> 4
			if (V[x] !=  V[y]) {
				PC += 2
			}
		case 0xA000:
			I = op & 0x0FFF

			PC += 2	
		case 0xB000:
			nnn := op & 0x0FFF
			PC = nnn + uint16(V[0])
		case 0xC000:
			rand.Seed(time.Now().UnixNano())
			x := (op & 0x0F00) >> 8
			kk := op & 0x00FF
			rand := rand.Intn(255)
			V[x] = byte(rand) & byte(kk)	
					
			PC += 2	

		case 0xD000:

 			x := (op & 0x0F00) >> 8
 			y := (op & 0x00F0) >> 4
 			n := (op & 0x000F)

			if(graphics.Draw(int(V[x]), int(V[y]), memory[I:I+n])) {
				V[0xF] = 0x01
			}

			graphics.Render()

			PC += 2	
		case 0xE000:
			x := (op & 0x0F00) >> 8
			switch(op & 0xF0FF) {
				case 0xE09E:
	            	key, _ := KeyboardReadByte()
					debug("\nKey pressed : 0x%04X", key)

	            	if key == V[x] {
	            		PC += 2
	            	}
	            	PC += 2

				case 0xE0A1:
	            	key, _ := KeyboardReadByte()
					debug("\nKey pressed : 0x%04X", key)

	            	if key != V[x] {
	            		PC += 2
	            	}
	            	PC += 2
				default:
					debug("Invalid opcode")

			} 	
		case 0xF000: 
	        x := (op & 0x0F00) >> 8
	        debug("\n%v", x)

	        switch(op & 0xF0FF) {
	            case 0xF007:
	            	// Fx07 - LD Vx, DT
	            	V[x] = DT
					
					PC += 2	
	            case 0xF00A:
	            	// Fx0A - LD Vx, K
	            	key, _ := KeyboardReadByte()
	            	V[x] = key

					debug("\nKey pressed : 0x%04X", key)

					PC += 2	

	            case 0xF015:
	            	// Fx15 - LD DT, Vx
					DT = V[x]
					PC += 2	

	            case 0xF018:
	            	// Fx18 - LD ST, Vx
					ST = V[x]
					PC += 2	

	            case 0xF01E:
	            	// Fx1E - ADD I, Vx
					I = I + uint16(V[x])
					PC += 2	

	            case 0xF029:
	            	// Fx29 - LD F, Vx
	            	I = uint16(V[x]) * uint16(0x05)

					PC += 2	

	            case 0xF033:
	            	// Fx33 - LD B, Vx
	            	memory[I] = V[x] / 100
	            	memory[I+1] = (V[x] / 10) % 10
	            	memory[I+2] = (V[x] % 100) % 10
					
					PC += 2	

	            case 0xF055:
	            	// Fx55 - LD [I], Vx
		            for i := uint16(0); byte(i) <= byte(x); i++ {
		            	memory[I+i] = V[i] 
		            }
					
					PC += 2	

	            case 0xF065:
	            	// Fx65 - LD Vx, [I]
		            for i := uint16(0); byte(i) <= byte(x); i++ {
		            	V[i] = memory[I+i]
		            }
					
					PC += 2	
	            default:
	            	debug("\nDefault")
	        }

		default:
			debug("\nWrong opcode, leaving the program")
			os.Exit(0)
	}
}

var keyMap = map[rune]byte{
	'1': 0x01, '2': 0x02, '3': 0x03, '4': 0x0C,
	'q': 0x04, 'w': 0x05, 'e': 0x06, 'r': 0x0D,
	'a': 0x07, 's': 0x08, 'd': 0x09, 'f': 0x0E,
	'z': 0x0A, 'x': 0x00, 'c': 0x0B, 'v': 0x0F,
}

func KeyboardReadByte() (byte, error) {
	debug("\nPress key !")
	event := termbox.PollEvent()

	debug("%e", event)

	key, ok := keyMap[event.Ch]
	if !ok {
		return 0x00, fmt.Errorf("unknown key: %v", event.Ch)
	}
	return key, nil
}

func debug(msg string, object ...interface{}) {
	if debugEnabled == true {
		fmt.Printf(msg, object...)
	}
}