package util_test

import (
	"testing"

	"github.com/lnobach/gonrg/util"
	"github.com/stretchr/testify/assert"
)

func TestEqualFoldASCIIOnly(t *testing.T) {
	assert.Equal(t, true, util.EqualFoldNonUnicode("foo.Example.org:8080", "foo.example.org:8080"))
	assert.Equal(t, false, util.EqualFoldNonUnicode("foo.example.org:8080", "bar.example.org:8080"))
	assert.Equal(t, false, util.EqualFoldNonUnicode("foo.exämple.org", "bar.EXÄMPLE.org"))
}
