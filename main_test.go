package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"apache-instruction-set-simulator/extras"
	"apache-instruction-set-simulator/machines"
	"apache-instruction-set-simulator/utils"
)

func Test_Apache8bits_Against_Programs(t *testing.T) {
	testCases := map[string]struct {
		programName, input, output string
		cycles                     int
	}{
		"sum.txt": {
			programName: "sum.txt",
			input:       "25\n" + "25\n",
			output:      "50\n",
			cycles:      999,
		},
		"sub.txt": {
			programName: "sub.txt",
			input:       "99\n" + "33\n",
			output:      "66\n",
			cycles:      999,
		},
		"square.txt": {
			programName: "square.txt",
			input:       "5\n",
			output:      "25\n",
			cycles:      999,
		},
		"fibonacci.txt": {
			programName: "fibonacci.txt",
			input:       "",
			output:      "1\n2\n3\n5\n8\n13\n21\n34\n55\n89\n144\n233\n",
			cycles:      44,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			in, err := utils.NewTestInput(testCase.input)
			if err != nil {
				t.Fatal(err)
			}
			defer in.Close()

			out := utils.NewTestOutput()

			var memory *extras.Memory16x8bits = extras.NewMemory16x8bits()
			memory.LoadProgram(testCase.programName)
			var machine *machines.Apache8bits = machines.NewApache8bits(memory, in, &out)
			machine.Run(testCase.cycles)

			assert.Equal(t, testCase.output, utils.ClearOutputForTesting(out.String()))
		})
	}
}
