package extras

import (
	"log"
	"strings"

	"apache-instruction-set-simulator/utils"
)

// RAM size (2 bits), 3 spaces
const memory3x8bitsSize uint8 = 0b11

type Memory3x8bits struct {
	MEMORY [memory3x8bitsSize]uint8
	SIZE   uint8
}

func (m *Memory3x8bits) Get(idx interface{}) interface{} {
	i := utils.CastInterfaceToUint8(idx)
	if i >= m.SIZE {
		log.Fatalf("Memory overflow, idx: %+v", idx)
	}
	return m.MEMORY[i]
}

func (m *Memory3x8bits) Set(idx interface{}, val interface{}) {
	i := utils.CastInterfaceToUint8(idx)
	if i >= m.SIZE {
		log.Fatalf("Memory overflow, idx: %+v", idx)
	}
	v := utils.CastInterfaceToUint8(val)
	m.MEMORY[i] = v
}

func (m *Memory3x8bits) Size() interface{} {
	return m.SIZE
}

func (m *Memory3x8bits) LoadProgram(program string) {
	pieces := strings.Split(program, "\n")
	for idx, piece := range pieces {
		if piece == "" {
			continue
		}
		m.Set(uint8(idx), utils.CastStringToUint8(piece, 2))
	}
}

func NewMemory3x8bits() *Memory3x8bits {
	device := &Memory3x8bits{}

	// set size
	device.SIZE = memory3x8bitsSize

	// RAM (2 bytes long)
	device.MEMORY = [memory3x8bitsSize]uint8{
		0b00000000,
		0b00000000,
		0b00000000,
	}

	return device
}
