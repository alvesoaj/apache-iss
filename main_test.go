package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"apache-instruction-set-simulator/extras"
	"apache-instruction-set-simulator/machines"
)

func TestApache8bits(t *testing.T) {
	testCases := map[string]struct {
		programName, input, output string
		cycles                     int
	}{
		"square.txt": {
			programName: "square.txt",
			input:       "5\n",
			output:      "25",
			cycles:      999,
		},
		"sum.txt": {
			programName: "sum.txt",
			input:       "25\n" + "25\n",
			output:      "50",
			cycles:      999,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			in, err := os.CreateTemp("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer in.Close()

			if _, err := io.WriteString(in, testCase.input); err != nil {
				t.Fatal(err)
			}

			if _, err := in.Seek(0, io.SeekStart); err != nil {
				t.Fatal(err)
			}

			var out bytes.Buffer

			var memory *extras.Memory16x8bits = extras.NewMemory16x8bits()
			memory.LoadProgram(testCase.programName)
			var machine *machines.Apache8bits = machines.NewApache8bits(memory, in, &out)
			machine.Run(testCase.cycles)

			assert.Equal(t, testCase.output, strings.TrimSpace(out.String()))
		})
	}
}
