package machines

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"apache-instruction-set-simulator/extras"
	"apache-instruction-set-simulator/utils"
)

func Test_Apache8bits(t *testing.T) {
	memory := extras.NewMemory3x8bits()
	assert.NotNil(t, memory)

	testCases := map[string]struct {
		program   string
		init      func(*Apache8bits)
		evaluator func(*testing.T, *Apache8bits, *os.File, bytes.Buffer)
	}{
		"LOAD R0": { // Load the ADDRESS into register 0
			init: func(machine *Apache8bits) {
				machine.REGISTERS[0] = 0b00000011
			},
			program: "0000 0001\n0000 1111", // register at 0 is 3, memory at 1 is 15
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(15), machine.REGISTERS[0])
			},
		},
		"STORE R0": { // Store content of register 0 into ADDRESS
			init: func(machine *Apache8bits) {
				machine.REGISTERS[0] = 0b00000011
			},
			program: "0001 0001\n0000 1111", // memory at 1 is 15, register at 0 is 3
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(3), machine.MEMORY.Get(uint8(1)))
			},
		},
		"JUMP R0 IF": { // Jump to line ADDRESS if register 0 is equal to 0
			init:    func(machine *Apache8bits) {},
			program: "0010 0010\n", // register at 0 is 0, pc is 0
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(2), machine.PC)
			},
		},
		"ADD R0": { // Add contents at ADDRESS to register 0
			init: func(machine *Apache8bits) {
				machine.REGISTERS[0] = 0b00000011
			},
			program: "0011 0001\n0000 1010", // memory at 1 is 10
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(13), machine.REGISTERS[0])
			},
		},
		"<<R0": { // Bitwise shift register 0 left
			init: func(machine *Apache8bits) {
				machine.REGISTERS[0] = 0b10000001
			},
			program: "0100 0000\n", // register at 0 is 129
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(2), machine.REGISTERS[0])
			},
		},
		"NOT R0": { // Bitwise NOT register 0
			init: func(machine *Apache8bits) {
				machine.REGISTERS[0] = 0b10010001
			},
			program: "0101 0000\n", // register at 0 is 145
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(110), machine.REGISTERS[0])
			},
		},
		"JUMP": { // Jump to line OPERAND
			init:    func(machine *Apache8bits) {},
			program: "0110 0010\n",
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(2), machine.PC)
			},
		},
		"STOP": { // Terminate the program (NOP)
			init:    func(machine *Apache8bits) {},
			program: "0111 0000\n",
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(1), machine.STOP)
			},
		},
		"LOAD R1": { // Load the ADDRESS into register 1
			init: func(machine *Apache8bits) {
				machine.REGISTERS[1] = 0b00000011
			},
			program: "1000 0001\n0000 1111", // register at 1 is 3, memory at 1 is 15
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(15), machine.REGISTERS[1])
			},
		},
		"STORE R1": { // Store content of register 1 into ADDRESS
			init: func(machine *Apache8bits) {
				machine.REGISTERS[1] = 0b00000011
			},
			program: "1001 0001\n0000 1111", // memory at 1 is 15, register at 1 is 3
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(3), machine.MEMORY.Get(uint8(1)))
			},
		},
		"JUMP R1 IF": { // Jump to line ADDRESS if register 1 is equal to 0
			init:    func(machine *Apache8bits) {},
			program: "1010 0010\n", // register at 1 is 0, pc is 0
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(2), machine.PC)
			},
		},
		"ADD R1": { // Add contents at ADDRESS to register 1
			init: func(machine *Apache8bits) {
				machine.REGISTERS[1] = 0b00000011
			},
			program: "1011 0001\n0000 1010", // memory at 1 is 10
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(13), machine.REGISTERS[1])
			},
		},
		"<<R1": { // Bitwise shift register 1 left
			init: func(machine *Apache8bits) {
				machine.REGISTERS[1] = 0b10000001
			},
			program: "1100 0000\n", // register at 0 is 129
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(2), machine.REGISTERS[1])
			},
		},
		"NOT R1": { // Bitwise NOT register 1
			init: func(machine *Apache8bits) {
				machine.REGISTERS[1] = 0b10010001
			},
			program: "1101 0000\n", // register at 1 is 145
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(110), machine.REGISTERS[1])
			},
		},
		"OUT R0": { // Outputs register 0
			init: func(machine *Apache8bits) {
				machine.REGISTERS[0] = 0b10010001
			},
			program: "1110 0000\n", // register at 0 is 145
			evaluator: func(t *testing.T, machine *Apache8bits, _ *os.File, out bytes.Buffer) {
				assert.Equal(t, "145\n", utils.ClearOutputForTesting(out.String()))
			},
		},
		"IN": { // Input into ADDRESS
			init:    func(machine *Apache8bits) {},
			program: "1111 0010\n",
			evaluator: func(t *testing.T, machine *Apache8bits, in *os.File, _ bytes.Buffer) {
				assert.Equal(t, uint8(100), machine.MEMORY.Get(uint8(2)))
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			memory.LoadProgram(testCase.program)

			in, err := utils.NewTestInput("100\n")
			if err != nil {
				t.Fatal(err)
			}
			defer in.Close()

			out := utils.NewTestOutput()

			machine := NewApache8bits(memory, in, &out)
			assert.NotNil(t, machine)

			testCase.init(machine)

			machine.Run(1)

			testCase.evaluator(t, machine, in, out)
		})
	}
}
