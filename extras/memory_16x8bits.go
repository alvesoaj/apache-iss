package extras

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"apache-instruction-set-simulator/utils"
)

// RAM size (5 bits), 16 spaces
const memory16x8bitsSize uint8 = 0b10000

type Memory16x8bits struct {
	MEMORY [memory16x8bitsSize]uint8
	SIZE   uint8
}

func (m *Memory16x8bits) Get(idx interface{}) interface{} {
	i := utils.CastInterfaceToUint8(idx)
	if i >= m.SIZE {
		log.Fatalf("Memory overflow, idx: %+v", idx)
	}
	return m.MEMORY[i]
}

func (m *Memory16x8bits) Set(idx interface{}, val interface{}) {
	i := utils.CastInterfaceToUint8(idx)
	if i >= m.SIZE {
		log.Fatalf("Memory overflow, idx: %+v", idx)
	}
	v := utils.CastInterfaceToUint8(val)
	m.MEMORY[i] = v
}

func (m *Memory16x8bits) Size() interface{} {
	return m.SIZE
}

func (m *Memory16x8bits) LoadProgram(programName string) {
	content, err := os.Open(fmt.Sprintf("./programs/%s", programName))
	if err != nil {
		log.Fatalf("File reading error: %+v", err)
	}

	var idx uint8 = 0
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		txt := scanner.Text()
		val := utils.CastStringToUint8(txt, 2)
		m.Set(idx, val)
		idx++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("File scanning error: %+v", err)
	}
}

func NewMemory16x8bits() *Memory16x8bits {
	device := &Memory16x8bits{}

	// set size
	device.SIZE = memory16x8bitsSize

	// RAM (16 bytes long)
	device.MEMORY = [memory16x8bitsSize]uint8{
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

	return device
}
