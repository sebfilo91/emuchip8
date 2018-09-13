package main

import(
	"fmt"
	"math/rand"
    "io/ioutil"
    "time"
    "os"
    "github.com/nsf/termbox-go"
    "graphics"
)

var opcode uint16

var memory [4096]byte

var stack[16]uint16

var SP uint16

var V [16]byte

var VF byte

var I uint16

// Program Counter
var PC uint16

// Delay Timer
var DT byte

// Sound Timer
var ST byte

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	loadRom("./roms/pong.rom")
	//loadRom("./roms/test.rom")

  	for i := 0; i < 50; i++ {
  		emulateCycle()
  	}
}

func loadRom(rom string) {
	// load at 0x200

	dat, err := ioutil.ReadFile(rom)
	if err != nil {
		fmt.Println(err)
	}

	loadAddress := int(0x200)
	PC = uint16(loadAddress)

	for i := 0; i < len(dat); i++ {
		memory[loadAddress + i] = dat[i]
		fmt.Printf("%#x ", memory[loadAddress + i] )
	}

	fmt.Printf("\nSTARTING MEMORY %#x \n", memory[0x200])
	fmt.Printf("STARTING PC %#x \n", PC)
	fmt.Printf("STARTING memory[PC] %#x \n", memory[PC])
}

func initialize() {

}

func emulateCycle() {
	opcode = uint16(memory[PC]) << 8 | uint16(memory[PC + 1])
	executeInstruction(opcode)
}

func executeInstruction(op uint16) {
	fmt.Printf("Op: %x\n", op)
	switch(op & 0xF000) {
		case 0x0000: 
			switch(op & 0x00FF) {
				case 0x00E0: 
					fmt.Println("CLS")
					//Clear graphics
					PC += 2
				case 0x00EE:
					fmt.Println("RET")
					PC = stack[SP]
					SP--
					PC += 2
			}
		case 0x1000:
			nnn := 0x0FFF & op
			PC = nnn
		case 0x2000:
			nnn := 0x0FFF & op
			SP++;
			stack[SP] = PC
			PC = nnn
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
					VF = carry
					V[x] = byte(result)

					PC += 2
				case 0x8005:
					result := uint16(V[x]) - uint16(V[y])

					VF = 0
					if V[x] > V[y] {
						VF = 1
					}

					V[x] = byte(result)

					PC += 2
				case 0x8006:

					PC += 2	

				case 0x8007:
					result := uint16(V[y]) - uint16(V[x])

					VF = 0
					if V[y] > V[x] {
						VF = 1
					}

					V[x] = byte(result)
					
					PC += 2

				case 0x800E:
					PC += 2	

			}

		case 0x9000:
			x := (op & 0x0F00) >> 8
			y := (op & 0x00F0) >> 4
			if (V[x] !=  V[y]) {
				PC += 2
			}
		case 0xA000:
			nnn := op & 0x0FFF
			I = nnn

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
			/*Dxyn - DRW Vx, Vy, nibble
Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.

The interpreter reads n bytes from memory, starting at the address stored in I. 
These bytes are then displayed as sprites on screen at coordinates (Vx, Vy). 
Sprites are XORed onto the existing screen. If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0.
 If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the screen.
 See instruction 8xy3 for more information on XOR, and section 2.4, Display, for more information on the Chip-8 screen and sprites.

			*/
 			x := (op & 0x0F00) >> 8
 			y := (op & 0x00F0) >> 4
 			n := (op & 0x000F)

 			for i := uint16(0); i < n; i++ {
 				Display(memory[i], V[x], V[y])	
 			}
					
			PC += 2	
		case 0xE000:
			x := op & 0x0F00
			switch(op & 0xF0FF) {
				case 0xE09E:
	            	key, _ := KeyboardReadByte()
					fmt.Println("Key pressed : %#x", key)

	            	if key == V[x] {
	            		PC += 2
	            	}
	            	PC += 2

				case 0xE0A1:
	            	key, _ := KeyboardReadByte()
					fmt.Println("Key pressed : %#x", key)

	            	if key != V[x] {
	            		PC += 2
	            	}
	            	PC += 2

			} 	
		case 0xF000: 
	        x := (op & 0x0F00) >> 8
	        fmt.Printf("%v", x)

	        switch(op & 0xF0FF) {
	            case 0xF007:
	            	// Fx07 - LD Vx, DT
	            	V[x] = DT
					
					PC += 2	
	            case 0xF00A:
	            	// Fx0A - LD Vx, K
	            	key, _ := KeyboardReadByte()
	            	V[x] = key

					fmt.Println("Key pressed : %#x", key)

					PC += 2	

	            case 0xF015:
	            	// Fx15 - LD DT, Vx
					DT = V[uint16(x)]
					PC += 2	

	            case 0xF018:
	            	// Fx18 - LD ST, Vx
					ST = V[uint16(x)]
					PC += 2	

	            case 0xF01E:
	            	// Fx1E - ADD I, Vx
					I = I + uint16(V[uint16(x)])
					PC += 2	

	            case 0xF029:
	            	// Fx29 - LD F, Vx
	            	//V[x]
	            	//I = 
					
					// c.I = uint16(c.V[x]) * uint16(0x05) why ?

					PC += 2	

	            case 0xF033:
	            	// Fx33 - LD B, Vx
	            	memory[I] = V[uint16(x)] / 100
	            	memory[I+1] = (V[uint16(x)] / 10) % 10
	            	memory[I+2] = (V[uint16(x)] % 100) % 10
					
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
	            	fmt.Println("Default")
	        }

		default:
			fmt.Println("Wrong opcode, leaving the program")
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
	fmt.Printf("Press key !")
	event := termbox.PollEvent()

	fmt.Println(event)

	key, ok := keyMap[event.Ch]
	if !ok {
		return 0x00, fmt.Errorf("unknown key: %v", event.Ch)
	}
	return key, nil
}

func Display(I byte, Vx byte, Vy byte) {

}

