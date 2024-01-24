package firstof_test

import (
	"http-everything/httpe/pkg/share/firstof"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	s := firstof.String("a", "b", "c")
	assert.Equal(t, "a", s)
}

func TestInt(t *testing.T) {
	i := firstof.Int(1, 2, 3)
	assert.Equal(t, 1, i)
}
