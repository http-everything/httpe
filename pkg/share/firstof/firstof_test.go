package firstof_test

import (
	"testing"

	"github.com/http-everything/httpe/pkg/share/firstof"

	"github.com/stretchr/testify/assert"
)

func TestFirstOfString(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		s := firstof.String("a", "b", "c")
		assert.Equal(t, "a", s)
	})
	t.Run("empty", func(t *testing.T) {
		s := firstof.String()
		assert.Equal(t, "", s)
	})
}

func TestFirstOfInt(t *testing.T) {
	t.Run("good", func(t *testing.T) {
		i := firstof.Int(1, 2, 3)
		assert.Equal(t, 1, i)
	})
	t.Run("empty", func(t *testing.T) {
		i := firstof.Int()
		assert.Equal(t, 0, i)
	})
	t.Run("zero", func(t *testing.T) {
		i := firstof.Int(0)
		assert.Equal(t, 0, i)
	})
}
