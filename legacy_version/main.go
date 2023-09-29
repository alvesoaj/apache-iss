package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// RAM size
const memorySize int = 16

// RAM (16 bytes long)
var memory [memorySize]uint8 = [memorySize]uint8{
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
	0b00000000,
}

// 2 General Purpose Registers (1 byte each)
var registers [2]uint8 = [2]uint8{
	0b00000000,
	0b00000000,
}

// Program Counter (It is [should be seen as] a 4 bits Special Purpose Register)
var pc uint8 = 0b0000

// Current Instruction Register (1 byte long)
var cir uint8 = 0b00000000

// Stop Register (1 bite long)
var stop uint8 = 0b0

// MASIC Instruction Set
// BINARY | OPCODE     | COMMENT
var instructions map[uint8]func(uint8) = map[uint8]func(uint8){
	// 0000   | LOAD R1    | Load the ADDRESS into register 1
	0b0000: func(a uint8) { registers[0] = memory[a] },
	// 0001   | STORE R1   | Store contents of register 1 into ADDRESS
	0b0001: func(a uint8) { memory[a] = registers[0] },
	// 0010   | JUMP R1 IF | Jump to line ADDRESS if register 1 is equal to 0
	0b0010: func(a uint8) {
		if registers[0] == 0b00000000 {
			pc = a
		}
	},
	// 0011   | ADD R1     | Add contents at ADDRESS to register 1
	0b0011: func(a uint8) { registers[0] += memory[a] },
	// 0100   | <<R1       | Bitwise shift register 1 left
	0b0100: func(a uint8) { registers[0] <<= 1 },
	// 0101   | NOT R1     | Bitwise NOT register 1
	0b0101: func(a uint8) { registers[0] = ^registers[0] },
	// 0110   | JUMP       | Jump to line OPERAND
	0b0110: func(a uint8) { pc = a },
	// 0111   | stop       | Terminate the program (NOP).
	0b0111: func(a uint8) { stop = 0b1 },
	// 1000   | LOAD R2    | Load the ADDRESS into register 2
	0b1000: func(a uint8) { registers[1] = memory[a] },
	// 1001   | STORE R2   | Store contents of register 2 into ADDRESS
	0b1001: func(a uint8) { memory[a] = registers[1] },
	// 1010   | JUMP R2 IF | Jump to line ADDRESS if register 2 is equal to 0
	0b1010: func(a uint8) {
		if registers[1] == 0b00000000 {
			pc = a
		}
	},
	// 1011   | ADD R2     | Add ADDRESS to register 2
	0b1011: func(a uint8) { registers[1] += memory[a] },
	// 1100   | <<R2       | Bitwise shift register 2 left
	0b1100: func(a uint8) { registers[1] <<= 1 },
	// 1101   | NOT R2     | Bitwise NOT register 2
	0b1101: func(a uint8) { registers[1] = ^registers[1] },
	// 1110   | OUT R1     | Outputs register 1
	0b1110: func(a uint8) { fmt.Println(registers[0]) },
	// 1111   | IN         | Input into ADDRESS
	0b1111: func(a uint8) {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		sVal, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Reading string error: %+v", err)
		}
		memory[a] = castStringToUint8(sVal, 10)
	},
}

var nonNumericRegex = regexp.MustCompile(`[^0-9]+`)

func run(cycles int) {
	for stop == 0b0 && cycles > 0 {
		cycles--
		// fetch
		cir = memory[pc]
		pc++
		// decode
		var instruction uint8 = cir >> 4
		var address uint8 = cir & 0b1111
		// execute
		instructions[instruction](address)
	}
}

func loadProgram(programName string) {
	content, err := os.Open(fmt.Sprintf("./programs/%s", programName))
	if err != nil {
		log.Fatalf("File reading error: %+v", err)
	}

	scanner := bufio.NewScanner(content)
	memoryCount := 0
	for scanner.Scan() {
		if memoryCount == memorySize {
			break
		}
		memory[memoryCount] = castStringToUint8(scanner.Text(), 2)
		memoryCount++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("File scanning error: %+v", err)
	}
}

func castStringToUint8(sVal string, base int) uint8 {
	sVal = nonNumericRegex.ReplaceAllString(sVal, "")
	nVal, err := strconv.ParseInt(sVal, base, 64)
	if err != nil {
		log.Fatalf("Parsing string to int error: %+v", err)
	}

	return uint8(nVal)
}

func main() {
	var programName string = os.Args[1]
	var cycles int = int(castStringToUint8(os.Args[2], 10))

	loadProgram(programName)
	run(cycles)

	fmt.Println("process finished")
}
