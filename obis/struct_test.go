package obis_test

import (
	"testing"

	"github.com/lnobach/gonrg/obis"
	"github.com/stretchr/testify/assert"
)

func TestPrettyValue(t *testing.T) {

	e := obis.OBISEntry{
		ExactKey:      "1-0:1.8.2*255",
		SimplifiedKey: "1.8.2",
		ValueNum:      12345678,
		ValueScale:    -4,
	}

	result := e.PrettyValue(true)

	assert.Equal(t, "1234.5678", result)

}

func TestPrettyValue_Negative(t *testing.T) {

	e := obis.OBISEntry{
		ExactKey:      "1-0:1.8.2*255",
		SimplifiedKey: "1.8.2",
		ValueNum:      -12345678,
		ValueScale:    -4,
	}

	result := e.PrettyValue(true)

	assert.Equal(t, "-1234.5678", result)

}

func TestPrettyValue_Text(t *testing.T) {

	e := obis.OBISEntry{
		ExactKey:      "1-0:1.8.2*255",
		SimplifiedKey: "1.8.2",
		ValueText:     "Foo",
	}

	result := e.PrettyValue(true)

	assert.Equal(t, "Foo", result)

}

func TestPrettyValue_Null(t *testing.T) {

	e := obis.OBISEntry{
		ExactKey:      "1-0:1.8.2*255",
		SimplifiedKey: "1.8.2",
	}

	result := e.PrettyValue(true)

	assert.Equal(t, "-", result)

}
