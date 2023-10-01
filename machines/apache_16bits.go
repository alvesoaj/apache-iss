package machines

import (
	"fmt"
	"io"
	"log"
	"os"

	"apache-instruction-set-simulator/extras"
	"apache-instruction-set-simulator/utils"
)

const apache16bitsMaxPCbits uint16 = 0b1111111111

type Apache16bits struct {
	REGISTERS    [4]uint16                     // 2 General Purpose Registers (1 word each)
	PC           uint16                        // Program Counter (1 word Special Purpose Register, max memory of 1024 spaces)
	CIR          uint16                        // Current Instruction Register (1 word long Special Purpose Register)
	STOP         uint8                         // Stop Register (1 bit [should be seen as a] long Special Purpose Register)
	INSTRUCTIONS map[uint8]func(uint8, uint16) // MASIC Instruction Set
	MEMORY       extras.Memory
}

// it will break the 16 bits in 3 pieces
// first 4 bits are for command, tue next 2 are for index 0 and the last 10 are for index 1
// cmd  idx0   idx1
// 0000 00     0000000000
func (m *Apache16bits) Run(cycles int) {
	for m.STOP == 0b0 && cycles > 0 {
		cycles--
		// fetch
		m.CIR = utils.CastInterfaceToUint16(m.MEMORY.Get(m.PC))
		m.PC++
		// decode
		var instruction uint8 = uint8(m.CIR >> 12)
		var addresses uint16 = m.CIR & 0b111111111111
		var address0 uint8 = uint8(addresses >> 10)
		var address1 uint16 = addresses & 0b1111111111
		// execute
		m.INSTRUCTIONS[instruction](address0, address1)
	}
}

func NewApache16bits(memory extras.Memory, in *os.File, out io.Writer) *Apache16bits {
	if in == nil {
		in = os.Stdin
	}

	if out == nil {
		out = os.Stdout
	}

	if !(utils.CastInterfaceToUint16(memory.Size()) <= apache16bitsMaxPCbits) {
		log.Fatalf("Memory is too big, max is: %d", apache16bitsMaxPCbits)
	}

	machine := &Apache16bits{
		MEMORY: memory,
	}

	// 4 General Purpose Registers
	machine.REGISTERS = [4]uint16{
		0b0000000000000000,
		0b0000000000000000,
		0b0000000000000000,
		0b0000000000000000,
	}

	// Program Counter
	machine.PC = 0b00000000

	// Current Instruction Register
	machine.CIR = 0b0000000000000000

	// Stop Register
	machine.STOP = 0b0

	//     BINARY | OPCODE      | COMMENT
	machine.INSTRUCTIONS = map[uint8]func(uint8, uint16){
		// 0000   | LOAD RX AX  | Load the ADDRESS X into register X
		0b0000: func(idx0 uint8, idx1 uint16) {
			machine.REGISTERS[idx0] = utils.CastInterfaceToUint16(machine.MEMORY.Get(idx1))
		},
		// 0001   | STORE RX AX | Store content of register X into ADDRESS X
		0b0001: func(idx0 uint8, idx1 uint16) { machine.MEMORY.Set(idx1, machine.REGISTERS[idx0]) },
		// 0010   | JUMP RX IF  | Jump to line ADDRESS X if register X is equal to 0
		0b0010: func(idx0 uint8, idx1 uint16) {
			if machine.REGISTERS[idx0] == 0b0000000000000000 {
				machine.PC = idx1
			}
		},
		// 0011   | ADD RX AX   | Add contents at ADDRESS X to register X
		0b0011: func(idx0 uint8, idx1 uint16) {
			machine.REGISTERS[idx0] += utils.CastInterfaceToUint16(machine.MEMORY.Get(idx1))
		},
		// 0100   | SUB RX AX   | Sub contents at ADDRESS X to register X
		0b0100: func(idx0 uint8, idx1 uint16) {
			machine.REGISTERS[idx0] -= utils.CastInterfaceToUint16(machine.MEMORY.Get(idx1))
		},
		// 0101   | MUT RX AX   | Mut contents at ADDRESS X to register X
		0b0101: func(idx0 uint8, idx1 uint16) {
			machine.REGISTERS[idx0] *= utils.CastInterfaceToUint16(machine.MEMORY.Get(idx1))
		},
		// 0110   | DIV RX AX   | Div contents at ADDRESS X to register X
		0b0110: func(idx0 uint8, idx1 uint16) {
			machine.REGISTERS[idx0] /= utils.CastInterfaceToUint16(machine.MEMORY.Get(idx1))
		},
		// 0111   | >>RX X      | Bitwise shift register X left, X times
		0b0111: func(idx0 uint8, idx1 uint16) { machine.REGISTERS[idx0] >>= uint16(idx1) },
		// 1000   | <<RX X      | Bitwise shift register X left, X times
		0b1000: func(idx0 uint8, idx1 uint16) { machine.REGISTERS[idx0] <<= uint16(idx1) },
		// 1001   | NOT RX      | Bitwise NOT register X
		0b1001: func(idx0 uint8, _ uint16) { machine.REGISTERS[idx0] = ^machine.REGISTERS[idx0] },
		// 1010   | JUMP        | Jump to line OPERAND
		0b1010: func(_ uint8, idx1 uint16) { machine.PC = idx1 },
		// 1010   |             |
		0b1011: func(_ uint8, _ uint16) {},
		// 1010   |             |
		0b1100: func(_ uint8, _ uint16) {},
		// 1011   | STOP        | Terminate the program (NOP)
		0b1101: func(_ uint8, _ uint16) { machine.STOP = 0b1 },
		// 1110   | OUT RX      | Outputs register X
		0b1110: func(idx0 uint8, _ uint16) { fmt.Fprintf(out, "%d\n", machine.REGISTERS[idx0]) },
		// 1111   | IN AX       | Input into ADDRESS
		0b1111: func(_ uint8, idx1 uint16) {
			var sVal string
			fmt.Fprint(out, "> ")
			fmt.Fscanf(in, "%s", &sVal)
			machine.MEMORY.Set(idx1, utils.CastStringToUint8(sVal, 10))
		},
	}

	return machine
}
