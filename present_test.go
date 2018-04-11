package natural

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPresentSV(t *testing.T) {
	expected := map[string]int64{
		// expected output, input
		"fem":                                                       5,
		"tjugotre":                                                  23,
		"nittionio":                                                 99,
		"etthundrasjuttiosex":                                       176,
		"tvåhundra":                                                 200,
		"ettusen":                                                   1000,
		"ettusenetthundra":                                          1100,
		"ettusenniohundratjugo":                                     1920,
		"två tusen":                                                 2000,
		"tre tusen":                                                 3000,
		"nittontusen åttahundrasextio":                              19860,
		"åttioniotusen":                                             89000,
		"femhundraåttioniotusen":                                    589000,
		"enmiljontvåhundratrettiofyratusen femhundrasextiosju":      1234567,
		"nio miljoner åttahundrasjuttiosextusen femhundrafyrtiotre": 9876543,
		"tolv miljoner":                                                          12000000,
		"tolv miljoner sexhundratusen":                                           12600000,
		"tvåhundrasextio miljoner":                                               260000000,
		"etthundratjugotre miljoner fyrahundrafemtiosextusen sjuhundraåttionio":  123456789,
		"niohundraåttiosju miljoner sexhundrafemtiofyratusen trehundratjugoett":  987654321,
		"etthundraelva miljoner etthundraelvatusen etthundraelva":                111111111,
		"niohundranittionio miljoner niohundranittioniotusen niohundranittionio": 999999999,
		"åtta miljarder":       8000000000,
		"åttahundra miljarder": 800000000000,
	}
	for s, i := range expected {
		assert.Equal(t, s, PresentSvSE(i))
	}
}

func TestPresentEN(t *testing.T) {
	expected := map[string]int64{
		// expected output, input
		"seventeen":                                         17,
		"ninety-seven":                                      97,
		"three hundred and forty-two":                       342,
		"nine hundred":                                      900,
		"nine hundred and fifty":                            950,
		"one thousand":                                      1000,
		"one thousand eleven":                               1011,
		"two thousand three hundred and forty-five":         2345,
		"ninety-nine thousand eight hundred and twenty-one": 99821,
	}
	for s, i := range expected {
		assert.Equal(t, s, PresentEnUS(i))
	}
}
