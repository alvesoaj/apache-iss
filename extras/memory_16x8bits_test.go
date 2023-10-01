package extras

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Memory16x8bits(t *testing.T) {
	memory := NewMemory16x8bits()
	assert.NotNil(t, memory)

	memory.LoadProgram("test.txt")

	assert.Equal(t, uint8(1), memory.Get(uint8(0)))
	assert.Equal(t, uint8(2), memory.Get(uint8(1)))
	assert.Equal(t, uint8(4), memory.Get(uint8(2)))
	assert.Equal(t, uint8(0), memory.Get(uint8(3)))

	memory.Set(uint8(3), uint8(8))
	assert.Equal(t, uint8(8), memory.Get(uint8(3)))

	assert.Equal(t, memory.SIZE, memory.Size())
}
