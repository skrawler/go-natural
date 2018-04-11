package natural

import (
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	multiplierMap = map[string]int{
		// eng
		"hundred":  100,
		"thousand": 1000,
		"million":  1000000,
		// swe
		// https://sv.wikipedia.org/wiki/Miljard
		// https://sv.wikipedia.org/wiki/Biljon
		// https://sv.wikipedia.org/wiki/Biljard_%28tal%29
		"hundra":    100,
		"tusen":     1000,
		"miljon":    1000000,
		"miljoner":  1000000,
		"miljard":   1000000000,
		"miljarder": 1000000000,
		"biljon":    1000000000000,
		"biljoner":  1000000000000,
		"biljard":   1000000000000000,
		"biljarder": 1000000000000000,
	}
)

// stringAsInt parses a string as a number
func prefixedStringAsInt(s string) (int64, error) {
	s = strings.TrimSpace(s)
	base := 10
	if strings.Contains(s, "0x") {
		// hex
		base = 16
		s = s[2:]
	} else if strings.Contains(s, "0b") {
		// binary
		base = 2
		s = s[2:]
	} else if strings.Contains(s, "0o") {
		// octal
		base = 8
		s = s[2:]
	} else {
		return 0, fmt.Errorf("no prefix")
	}
	return strconv.ParseInt(s, base, 64)
}

// NumberStringToBig converts a string representation of a number such as "123" to a Decimal representation
func NumberStringToBig(s string) (decimal.Decimal, error) {
	s = strings.TrimSpace(s)
	n, err := prefixedStringAsInt(s)
	if err == nil {
		return decimal.NewFromFloat(float64(n)), nil
	}
	return decimal.NewFromString(s)
}

// ParseNumber parses cardinal numbers (like "five") in English and Swedish
func ParseNumber(s string) (decimal.Decimal, error) {
	if num, err := NumberStringToBig(s); err == nil {
		return num, nil
	}
	if res, err := ParseNumberSwedish(s); err == nil {
		return res, nil
	}
	if res, err := ParseNumberEnglish(s); err == nil {
		return res, nil
	}
	return decimal.NewFromFloat(0), fmt.Errorf("Cannot parse number '%s'", s)
}

func isNumericString(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

var (
	multiplierSvSERegex = regexp.MustCompile(`^(?P<num>[\d]+) (?P<size>hundra|tusen|miljon(er)?|miljard(er)?|biljon(er)?|biljard(er)?)+$`)
	wholeAndFraction    = regexp.MustCompile(`^(?P<arg1>.*) och (?P<arg2>.*)$`)
	wholeCommaDecimal   = regexp.MustCompile(`^(?P<arg1>.*) komma (?P<arg2>.*)$`)
	scaleDataSV         = []struct {
		singular string
		plural   string
		scale    float64
	}{
		{
			// 1,000,000,000,000 - 999,999,999,999,999 (1 biljard)
			"biljard",
			"biljarder",
			1000000000000000.,
		},
		{
			// 1,000,000,000,000 - 999,999,999,999,999 (1 biljon)
			"biljon",
			"biljoner",
			1000000000000.,
		},
		{
			// 1,000,000,000 - 999,999,999,999 (1 miljard)
			"miljard",
			"miljarder",
			1000000000.,
		},
		{
			// 1,000,000 - 999,999,999 (1 miljon)
			"miljon",
			"miljoner",
			1000000.,
		},
	}
)

// ParseNumberSwedish parses a natural number in written Swedish
func ParseNumberSwedish(s string) (decimal.Decimal, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	var err error
	res := decimal.NewFromFloat(0)
	if s == "" {
		return res, nil
	}
	if s == "hälften" {
		return decimal.NewFromString("0.5")
	}

	match := multiplierSvSERegex.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		res, err = mapToMultiplier(match[0][1], match[0][2])
		return res, err
	}

	// https://sv.wikipedia.org/wiki/Namn_p%C3%A5_stora_tal
	s = strings.Replace(s, " hundra", "hundra", -1)
	s = strings.Replace(s, " tusen", "tusen", -1)
	s = strings.Replace(s, " miljon", "miljon", -1)   // 10^6
	s = strings.Replace(s, " miljard", "miljard", -1) // 10^9
	s = strings.Replace(s, " biljon", "biljon", -1)   // 10^12
	s = strings.Replace(s, " biljard", "biljard", -1) // 10^15

	// "fem komma två"
	match = wholeCommaDecimal.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		res, err = ParseNumberSwedish(match[0][1])
		if err != nil {
			return res, err
		}
		var dec decimal.Decimal
		dec, err = ParseNumberSwedish(match[0][2])
		if err != nil {
			return res, err
		}
		//ten := decimal.NewFromFloat(10.)
		// XXX the Pow() in use is not mainline yet: https://github.com/shopspring/decimal/pull/29
		//scale := ten.Pow(decimal.NewFromFloat(float64(len(dec.String()))))
		scale := decimal.NewFromFloat(math.Pow(10, float64(len(dec.String()))))
		return dec.Div(scale).
			Add(res), nil
	}

	// "femton och tre fjärdedelar" = 15.75
	match = wholeAndFraction.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		var frac decimal.Decimal
		res, err = ParseNumberSwedish(match[0][1])
		if err != nil {
			return res, err
		}
		frac, err = parseFractions(match[0][2])
		if err != nil {
			return res, err
		}
		return res.Add(frac), nil
	}

	for _, d := range scaleDataSV {
		if strings.Contains(s, d.singular) {
			for i := int64(1); i <= 999; i++ {
				prefix := ""
				if i == 1 {
					prefix = "en" + d.singular
				} else {
					prefix = PresentSvSE(i) + d.plural
				}
				if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
					res, err = ParseNumberSwedish(s[len(prefix):])
					return decimal.NewFromFloat(float64(i)).
						Mul(decimal.NewFromFloat(d.scale)).
						Add(res), err
				}
			}
		}
	}

	// 1,000 - 999,999
	for i := int64(1); i <= 999; i++ {
		if i == 1 && len(s) >= 5 && s[0:5] == "tusen" {
			res, err = ParseNumberSwedish(s[5:])
			return decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(1000.)).
				Add(res), err
		}
		if i == 1 && len(s) >= 7 && s[0:7] == "ettusen" {
			res, err = ParseNumberSwedish(s[7:])
			return decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(1000.)).
				Add(res), err
		}
		if i == 1 && len(s) >= 11 && s[0:11] == "hundratusen" {
			res, err = ParseNumberSwedish(s[11:])
			return decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(100000.)).
				Add(res), err
		}
		if i == 1 && len(s) >= 14 && s[0:14] == "'etthundratusen" {
			res, err = ParseNumberSwedish(s[14:])
			return decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(100000.)).
				Add(res), err
		}
		prefix := PresentSvSE(i) + "tusen"
		if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
			res, err = ParseNumberSwedish(s[len(prefix):])
			return decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(1000.)).
				Add(res), err
		}
	}

	// 100 - 1999 ("nittonhundra"... ej "ettusenniohundra")
	for prefix, i := range numbersToTwentySvSE {
		prefix = prefix + "hundra"
		if i == 1 && len(s) >= 6 && s[0:6] == "hundra" {
			res, err = ParseNumberSwedish(s[6:])
			if err != nil {
				return res, err
			}
			res = decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(100.)).
				Add(res)
			return res, nil
		}
		if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
			res, err = ParseNumberSwedish(s[len(prefix):])
			if err != nil {
				return res, err
			}
			res = decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(100.)).
				Add(res)
			return res, nil
		}
	}

	// 20 - 100
	for _, prefix := range tensSvSE {
		if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
			if tens, _err := arrayIndex(prefix, tensSvSE); _err == nil {
				res, err = ParseNumberSwedish(s[len(prefix):])
				if err != nil {
					return res, err
				}
				_tens := decimal.NewFromFloat(float64(tens * 10))
				return _tens.Add(res), nil
			}
		}
	}

	// 1 - 20
	if v, ok := numbersToTwentySvSE[s]; ok {
		return decimal.NewFromFloat(float64(v)), nil
	}
	if s == "en" {
		return decimal.NewFromFloat(1.), nil
	}

	// "tre fjärdedelar"
	res, err = parseFractions(s)
	return res, err
}

// ParseNumberEnglish parses a natural number in written English
func ParseNumberEnglish(s string) (decimal.Decimal, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	// XXX handle higher numbers

	if s == "" {
		return decimal.NewFromFloat(0), nil
	}

	// https://en.wikipedia.org/wiki/Names_of_large_numbers
	s = strings.Replace(s, " hundred", "hundred", -1)
	s = strings.Replace(s, "hundred ", "hundred", -1)
	s = strings.Replace(s, " thousand", "thousand", -1)
	s = strings.Replace(s, "thousand ", "thousand", -1)
	s = strings.Replace(s, " million", "million", -1) // 10^6
	s = strings.Replace(s, "million ", "million", -1)
	s = strings.Replace(s, " billion", "billion", -1) // 10^9
	s = strings.Replace(s, "billion ", "billion", -1)
	s = strings.Replace(s, " trillion", "trillion", -1) // 10^12
	s = strings.Replace(s, "trillion", "trillion", -1)

	// 100 - 1999
	for prefix, i := range numbersToTwentyEnUS {
		prefix = prefix + "hundred"
		if i == 1 && len(s) >= 7 && s[0:7] == "hundred" {
			res, err := ParseNumberEnglish(s[7:])
			if err != nil {
				return res, err
			}
			res = decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(100.)).
				Add(res)
			return res, nil
		}
		if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
			res, err := ParseNumberEnglish(s[len(prefix):])
			if err != nil {
				return res, err
			}
			res = decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(100.)).
				Add(res)
			return res, nil
		}
	}

	// 20 - 100
	for _, prefix := range tensEnUS {
		if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
			if tens, _err := arrayIndex(prefix, tensEnUS); _err == nil {
				res, err := ParseNumberEnglish(s[len(prefix):])
				if err != nil {
					return res, err
				}
				_tens := decimal.NewFromFloat(float64(tens * 10))
				return _tens.Add(res), nil
			}
		}
	}

	// 1 - 20
	if v, ok := numbersToTwentyEnUS[s]; ok {
		return decimal.NewFromFloat(float64(v)), nil
	}

	return decimal.NewFromFloat(0), fmt.Errorf("error")
}

func mapToMultiplier(num, multiplier string) (decimal.Decimal, error) {
	var res decimal.Decimal
	if _, ok := multiplierMap[multiplier]; !ok {
		return res, fmt.Errorf("key not found")
	}
	n, err := decimal.NewFromString(num)
	if err != nil {
		return res, err
	}
	res = n.Mul(decimal.NewFromFloat(float64(multiplierMap[multiplier])))
	return res, err
}

func getFraction(s string) *big.Rat {
	// s: tredjedelars => tredjedelar, tredjedels => tredjedel
	if len(s) > 1 && s[len(s)-1:] == "s" {
		s = s[0 : len(s)-1]
	}
	// tredjedelar => tredjedel
	if len(s) > 2 && s[len(s)-2:] == "ar" {
		s = s[0 : len(s)-2]
	}
	if v, ok := fractionsSvSE[s]; ok {
		r := new(big.Rat)
		r.SetString(v)
		return r
	}
	fmt.Println("ERROR getFraction failed to parse", s)
	return nil
}

func parseFractions(s string) (decimal.Decimal, error) {
	var res decimal.Decimal
	x := strings.SplitN(s, " ", 2)

	if len(x) == 1 {
		return res, fmt.Errorf("parseFractions failed %s", s)
	}
	if len(x) == 2 {
		num, err := ParseNumber(x[0])
		if err != nil {
			return res, err
		}
		fraction := getFraction(x[1])
		if fraction == nil {
			return res, fmt.Errorf("getFraction failed %s", x[1])
		}
		// XXX hack, loss of precision:
		frac, err := decimal.NewFromString(fraction.FloatString(16))
		if err != nil {
			return res, err
		}
		return num.Mul(frac), nil
	}
	return res, fmt.Errorf("nothing parsed")
}
