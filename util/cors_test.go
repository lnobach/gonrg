package util_test

import (
	"testing"

	"github.com/lnobach/gonrg/util"
	"github.com/stretchr/testify/assert"
)

func TestDomainMatches(t *testing.T) {
	assert.Equal(t, true, util.DomainMatches("foo.example.org", "foo.example.org"))
	assert.Equal(t, true, util.DomainMatches("foo.example.org", "*"))
	assert.Equal(t, true, util.DomainMatches("foo.example.org", "*g"))
	assert.Equal(t, true, util.DomainMatches("org", "*org"))
	assert.Equal(t, true, util.DomainMatches("rg", "*g"))
	assert.Equal(t, true, util.DomainMatches("g", "*g"))
	assert.Equal(t, false, util.DomainMatches("foo.example.org", "*toolongpattern.example.org"))
	assert.Equal(t, true, util.DomainMatches("foo.example.org", "*.example.org"))
	assert.Equal(t, false, util.DomainMatches("", "*.example.org"))
	assert.Equal(t, false, util.DomainMatches("foo", ""))
	assert.Equal(t, true, util.DomainMatches("example.org", "*example.org"))
	assert.Equal(t, true, util.DomainMatches("example.org", "*Example.org"))
}
