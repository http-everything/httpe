package merge_test

import (
	"github.com/stretchr/testify/assert"
	"http-everything/httpe/pkg/share/merge"
	"testing"
)

func TestStringMapsI(t *testing.T) {
	map1 := map[string]string{
		"keY1": "a",
		"Key2": "b",
		"key3": "c",
	}
	map2 := map[string]string{
		"Key1": "a-new",
		"kEy3": "b-new",
		"key4": "d-new",
	}
	result := merge.StringMapsI(map1, map2)
	wants := map[string]string{
		"key1": "a",
		"key2": "b",
		"key3": "c",
		"key4": "d-new",
	}
	assert.Equal(t, wants, result)
}
