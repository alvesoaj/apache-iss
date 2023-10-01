package extras

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Memory3x8bits(t *testing.T) {
	memory := NewMemory3x8bits()
	assert.NotNil(t, memory)

	memory.LoadProgram("0000 0001\n0000 0010\n0000 0011")

	assert.Equal(t, uint8(1), memory.Get(uint8(0)))

	memory.Set(uint8(0), uint8(2))
	assert.Equal(t, uint8(2), memory.Get(uint8(0)))
}
