package machines

import (
	"fmt"
	"io"
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

func NewApache8bits(memory extras.Memory, in *os.File, out io.Writer) *Apache8bits {
	if in == nil {
		in = os.Stdin
	}

	if out == nil {
		out = os.Stdout
	}

	machine := &Apache8bits{
		MEMORY: memory,
	}

	// 2 General Purpose Registers (1 byte each)
	machine.REGISTERS = [2]uint8{
		0b00000000,
		0b00000000,
	}

	// Program Counter (It is [should be seen as] a 4 bits Special Purpose Register)
	machine.PC = 0b0000

	// Current Instruction Register (It is [should be seen as] a 1 byte long)
	machine.CIR = 0b00000000

	// Stop Register (It is [should be seen as] a 1 bite long)
	machine.STOP = 0b0

	// BINARY | OPCODE     | COMMENT
	machine.INSTRUCTIONS = map[uint8]func(uint8){
		// 0000   | LOAD R0    | Load the ADDRESS into register 0
		0b0000: func(idx uint8) { machine.REGISTERS[0] = utils.CastInterfaceToUint8(machine.MEMORY.Get(idx)) },
		// 0001   | STORE R0   | Store content of register 0 into ADDRESS
		0b0001: func(idx uint8) { machine.MEMORY.Set(idx, machine.REGISTERS[0]) },
		// 0010   | JUMP R0 IF | Jump to line ADDRESS if register 0 is equal to 0
		0b0010: func(idx uint8) {
			if machine.REGISTERS[0] == 0b00000000 {
				machine.PC = idx
			}
		},
		// 0011   | ADD R0     | Add contents at ADDRESS to register 0
		0b0011: func(idx uint8) { machine.REGISTERS[0] += utils.CastInterfaceToUint8(machine.MEMORY.Get(idx)) },
		// 0100   | <<R0       | Bitwise shift register 0 left
		0b0100: func(_ uint8) { machine.REGISTERS[0] <<= 1 },
		// 0101   | NOT R0     | Bitwise NOT register 0
		0b0101: func(_ uint8) { machine.REGISTERS[0] = ^machine.REGISTERS[0] },
		// 0110   | JUMP       | Jump to line OPERAND
		0b0110: func(idx uint8) { machine.PC = idx },
		// 0111   | STOP       | Terminate the program (NOP)
		0b0111: func(_ uint8) { machine.STOP = 0b1 },
		// 1000   | LOAD R1    | Load the ADDRESS into register 1
		0b1000: func(idx uint8) { machine.REGISTERS[1] = utils.CastInterfaceToUint8(machine.MEMORY.Get(idx)) },
		// 1001   | STORE R1   | Store contents of register 1 into ADDRESS
		0b1001: func(idx uint8) { machine.MEMORY.Set(idx, machine.REGISTERS[1]) },
		// 1010   | JUMP R1 IF | Jump to line ADDRESS if register 1 is equal to 0
		0b1010: func(idx uint8) {
			if machine.REGISTERS[1] == 0b00000000 {
				machine.PC = idx
			}
		},
		// 1011   | ADD R1     | Add ADDRESS to register 1
		0b1011: func(idx uint8) { machine.REGISTERS[1] += utils.CastInterfaceToUint8(machine.MEMORY.Get(idx)) },
		// 1100   | <<R1       | Bitwise shift register 1 left
		0b1100: func(_ uint8) { machine.REGISTERS[1] <<= 1 },
		// 1101   | NOT R1     | Bitwise NOT register 1
		0b1101: func(_ uint8) { machine.REGISTERS[1] = ^machine.REGISTERS[1] },
		// 1110   | OUT R0     | Outputs register 0
		0b1110: func(_ uint8) { fmt.Fprintf(out, "%d\n", machine.REGISTERS[0]) },
		// 1111   | IN         | Input into ADDRESS
		0b1111: func(idx uint8) {
			var sVal string
			fmt.Fprint(out, "> ")
			fmt.Fscanf(in, "%s", &sVal)
			machine.MEMORY.Set(idx, utils.CastStringToUint8(sVal, 10))
		},
	}

	return machine
}
