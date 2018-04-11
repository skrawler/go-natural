package natural

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestParseNumber(t *testing.T) {
	expected := map[string]string{
		"":       "0",
		"2":      "2",
		"59":     "59",
		"5.2":    "5.2",
		"0x20":   "32", // hex
		"0b1001": "9",  // binary
		// eng
		"sixteen":              "16",
		"fiftythree":           "53",
		"sixty":                "60",
		"onehundredsixty":      "160",
		"ninehundredfiftytwo":  "952",
		"ninehundred fiftytwo": "952",
		//"three thousand ninetyfive": "3095",
		//"eighthundredthousandnine":  "800009", // XXX currently returns 1809
		// swe
		"tretton":                                             "13",
		"tjugo":                                               "20",
		"femtionio":                                           "59",
		"sextio":                                              "60",
		"hundranittionio":                                     "199",
		"etthundranittionio":                                  "199",
		"två hundra":                                          "200",
		"2 hundra":                                            "200",
		"sexhundrasextiosex":                                  "666",
		"niohundranittionio":                                  "999",
		"nittonhundratrettiotvå":                              "1932",
		"tusen":                                               "1000",
		"ettusen":                                             "1000",
		"ettusenniohundratrettiotvå":                          "1932",
		"2 tusen":                                             "2000",
		"tvåtusen":                                            "2000",
		"tvåtusenelva":                                        "2011",
		"tvåtusenetthundranio":                                "2109",
		"tvåtusenetthundranittionio":                          "2199",
		"niotusen":                                            "9000",
		"nio tusen":                                           "9000",
		"9999":                                                "9999",
		"åttioniotusen":                                       "89000",
		"nittioniotusenniohundranittioåtta":                   "99998",
		"hundratusen":                                         "100000",
		"hundra tusen":                                        "100000",
		"etthundratusen":                                      "100000",
		"etthundra tusen":                                     "100000",
		"100 tusen":                                           "100000",
		"tvåhundratusentretusennitton":                        "203019",
		"enmiljontvåhundratrettiofyratusenfemhundrasextiosju": "1234567",
		"åtta miljoner":                                       "8000000",
		"8 miljoner":                                          "8000000",
		"etthundraelvamiljoneretthundraelvatusenetthundraelva":                "111111111",
		"niohundranittioniomiljonerniohundranittioniotusenniohundranittionio": "999999999",
		"åttamiljarder":  "8000000000",
		"åtta miljarder": "8000000000",
		"8 miljarder":    "8000000000",
		"åttabiljoner":   "8000000000000",
		"åtta biljoner":  "8000000000000",
		"8 biljoner":     "8000000000000",
		"1 biljard":      "1000000000000000",
		"åtta biljarder": "8000000000000000",

		// swe - fractions
		"fem komma två":              "5.2",
		"hälften":                    "0.5",
		"en halv":                    "0.5",
		"en femtedel":                "0.2",
		"tre fjärdedelar":            "0.75",
		"tre femtedel":               "0.6",
		"tre femtedels":              "0.6",
		"tre femtedelar":             "0.6",
		"tre femtedelars":            "0.6",
		"femton och tre fjärdedelar": "15.75",
		"en sjundedel":               "0.1428571428571429",
	}
	for s, i := range expected {
		m, err := ParseNumber(s)
		assert.Equal(t, nil, err)
		expectedVal, err := decimal.NewFromString(i)
		assert.Equal(t, nil, err)
		assert.Equal(t, expectedVal.String(), m.String(), "input: "+s)
	}
}
