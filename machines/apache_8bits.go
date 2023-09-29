package machines

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"apache-instruction-set-simulator/extras"
	"apache-instruction-set-simulator/utils"
)

type Apache8bits struct {
	REGISTERS    [2]uint8              // 2 General Purpose Registers (1 byte each)
	PC           uint8                 // Program Counter (It is [should be seen as] a 4 bits Special Purpose Register)
	CIR          uint8                 // Current Instruction Register (1 byte long Special Purpose Register)
	STOP         uint8                 // Stop Register (1 bite long Special Purpose Register)
	INSTRUCTIONS map[uint8]func(uint8) // MASIC Instruction Set
	MEMORY       extras.Memory
}

func (m *Apache8bits) Run(cycles int) {
	for m.STOP == 0b0 && cycles > 0 {
		cycles--
		// fetch
		m.CIR = utils.CastInterfaceToUint8(m.MEMORY.Get(m.PC))
		m.PC++
		// decode
		var instruction uint8 = m.CIR >> 4
		var address uint8 = m.CIR & 0b1111
		// execute
		m.INSTRUCTIONS[instruction](address)
	}
}

func NewApache8bits(memory extras.Memory) Apache8bits {
	machine := Apache8bits{
		MEMORY: memory,
	}

	// 2 General Purpose Registers (1 byte each)
	machine.REGISTERS = [2]uint8{
		0b00000000,
		0b00000000,
	}

	// Program Counter (It is [should be seen as] a 4 bits Special Purpose Register)
	machine.PC = 0b0000

	// Current Instruction Register (1 byte long)
	machine.CIR = 0b00000000

	// Stop Register (1 bite long)
	machine.STOP = 0b0

	// BINARY | OPCODE     | COMMENT
	machine.INSTRUCTIONS = map[uint8]func(uint8){
		// 0000   | LOAD R1    | Load the ADDRESS into register 1
		0b0000: func(idx uint8) { machine.REGISTERS[0] = utils.CastInterfaceToUint8(machine.MEMORY.Get(idx)) },
		// 0001   | STORE R1   | Store contents of register 1 into ADDRESS
		0b0001: func(idx uint8) { machine.MEMORY.Set(idx, machine.REGISTERS[0]) },
		// 0010   | JUMP R1 IF | Jump to line ADDRESS if register 1 is equal to 0
		0b0010: func(idx uint8) {
			if machine.REGISTERS[0] == 0b00000000 {
				machine.PC = idx
			}
		},
		// 0011   | ADD R1     | Add contents at ADDRESS to register 1
		0b0011: func(idx uint8) { machine.REGISTERS[0] += utils.CastInterfaceToUint8(machine.MEMORY.Get(idx)) },
		// 0100   | <<R1       | Bitshift register 1 left
		0b0100: func(idx uint8) { machine.REGISTERS[0] <<= 1 },
		// 0101   | NOT R1     | Bitwise NOT register 1
		0b0101: func(idx uint8) { machine.REGISTERS[0] = ^machine.REGISTERS[0] },
		// 0110   | JUMP       | Jump to line OPERAND
		0b0110: func(idx uint8) { machine.PC = idx },
		// 0111   | stop       | Terminate the program (NOP).
		0b0111: func(idx uint8) { machine.STOP = 0b1 },
		// 1000   | LOAD R2    | Load the ADDRESS into register 2
		0b1000: func(idx uint8) { machine.REGISTERS[1] = utils.CastInterfaceToUint8(machine.MEMORY.Get(idx)) },
		// 1001   | STORE R2   | Store contents of register 2 into ADDRESS
		0b1001: func(idx uint8) { machine.MEMORY.Set(idx, machine.REGISTERS[1]) },
		// 1010   | JUMP R2 IF | Jump to line ADDRESS if register 2 is equal to 0
		0b1010: func(idx uint8) {
			if machine.REGISTERS[1] == 0b00000000 {
				machine.PC = idx
			}
		},
		// 1011   | ADD R2     | Add ADDRESS to register 2
		0b1011: func(idx uint8) { machine.REGISTERS[1] += utils.CastInterfaceToUint8(machine.MEMORY.Get(idx)) },
		// 1100   | <<R2       | Bitshift register 2 left
		0b1100: func(idx uint8) { machine.REGISTERS[1] <<= 1 },
		// 1101   | NOT R2     | Bitwise NOT register 2
		0b1101: func(idx uint8) { machine.REGISTERS[1] = ^machine.REGISTERS[1] },
		// 1110   | OUT R1     | Outputs register 1
		0b1110: func(idx uint8) { fmt.Println(machine.REGISTERS[0]) },
		// 1111   |            |
		0b1111: func(idx uint8) {
			fmt.Print("> ")
			reader := bufio.NewReader(os.Stdin)
			sVal, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Reading string error: %+v", err)
			}
			machine.MEMORY.Set(idx, utils.CastStringToUint8(sVal, 10))
		},
	}

	return machine
}