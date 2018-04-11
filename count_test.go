package natural

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestParseCount(t *testing.T) {
	m, err := ParseCount("")
	assert.NotEqual(t, nil, err)

	expected := map[string]string{
		"13": "13",
		// eng
		"38:th": "38",
		"33:rd": "33",
		"fifth": "5",
		// swe
		"trettonde":             "13",
		"tjugotredje":           "23",
		"hundranittonde":        "119",
		"etthundranittiosjunde": "197",
		"nittonhundranionde":    "1909",
		"13:e":                  "13",
		"229:a":                 "229",
	}
	for s, i := range expected {
		m, err = ParseCount(s)
		assert.Equal(t, nil, err)

		expectedVal, err := decimal.NewFromString(i)
		assert.Equal(t, nil, err)
		assert.Equal(t, expectedVal.String(), m.String(), "input: "+s)
	}
}

func TestPresentCount(t *testing.T) {
	expectedSV := map[int]string{
		9:  "nionde",
		13: "trettonde",
		91: "nittioförsta",
	}
	for n, expect := range expectedSV {
		assert.Equal(t, expect, PresentCountSwedish(n))
	}

	expectedEN := map[int]string{
		9:  "ninth",
		13: "thirteenth",
		91: "ninetyfirst",
	}
	for n, expect := range expectedEN {
		assert.Equal(t, expect, PresentCountEnglish(n))
	}
}

func TestPresentCountShort(t *testing.T) {
	expectedSV := map[int]string{
		9:  "9:e",
		13: "13:e",
		91: "91:a",
	}
	for n, expect := range expectedSV {
		assert.Equal(t, expect, PresentCountShortSwedish(n))
	}

	expectedEN := map[int]string{
		3:  "3:rd",
		9:  "9:th",
		13: "13:th",
		91: "91:st",
	}
	for n, expect := range expectedEN {
		assert.Equal(t, expect, PresentCountShortEnglish(n))
	}
}
