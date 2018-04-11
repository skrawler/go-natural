package natural

import (
	"fmt"
	"log"
	"strings"

	"github.com/shopspring/decimal"
)

// PresentCountEnglish renders the count in Swedish
func PresentCountEnglish(n int) string {

	// 20 - 100
	if n > 20 && n < 100 {
		tens := (n / 10) % 10
		ones := n % 10
		return tensEnUS[tens] + countNamesEnUS[ones]
	}

	// 1 - 20
	if n < len(countNamesEnUS) {
		return countNamesEnUS[n]
	}

	return "XXX-PresentCountEnglish"
}

// PresentCountShortSwedish renders a short count in Swedish, such as "13:e"
func PresentCountShortSwedish(n int) string {
	s := PresentCountSwedish(n)
	last := s[len(s)-1]
	if last == 'a' {
		return fmt.Sprintf("%d:a", n)
	}
	return fmt.Sprintf("%d:e", n)
}

// PresentCountShortEnglish renders a short count in English, such as "13:th"
func PresentCountShortEnglish(n int) string {
	s := PresentCountEnglish(n)
	if len(s) < 3 {
		log.Fatal("invalid input", n)
	}
	last := s[len(s)-2 : len(s)]
	return fmt.Sprintf("%d:%s", n, last)
}

// PresentCountSwedish renders the count in Swedish, such as "trettonde"
func PresentCountSwedish(n int) string {

	// 20 - 100
	if n > 20 && n < 100 {
		tens := (n / 10) % 10
		ones := n % 10
		return tensSvSE[tens] + countNamesSvSE[ones]
	}

	// 1 - 20
	if n < len(countNamesSvSE) {
		return countNamesSvSE[n]
	}

	return "XXX-PresentCountSwedish"
}

// ParseCount parses ordinal numbers (like "fifth") in English and Swedish
func ParseCount(s string) (decimal.Decimal, error) {
	num, err := decimal.NewFromString(s)
	if err == nil {
		return num, nil
	}
	s = strings.ToLower(s)
	if res, err := parseCountSwedish(s); err == nil {
		return res, nil
	}
	if res, err := parseCountEnglish(s); err == nil {
		return res, nil
	}
	return decimal.NewFromFloat(0), fmt.Errorf("count: parse error: %s", s)
}

func parseCountSwedish(s string) (decimal.Decimal, error) {
	var res decimal.Decimal

	// 100 - 1999 ("nittonhundra"... ej "ettusenniohundra")
	for prefix, i := range numbersToTwentySvSE {
		prefix = prefix + "hundra"
		if i == 1 && len(s) >= 6 && s[0:6] == "hundra" {
			res, err := parseCountSwedish(s[6:])
			if err != nil {
				return res, err
			}
			res = decimal.NewFromFloat(float64(i)).
				Mul(decimal.NewFromFloat(100.)).
				Add(res)
			return res, nil
		}
		if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
			res, err := parseCountSwedish(s[len(prefix):])
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
			if tens, err := arrayIndex(prefix, tensSvSE); err == nil {
				restStr := s[len(prefix):]
				if restStr != "" {
					res, err = parseCountSwedish(restStr)
					if err != nil {
						return res, err
					}
				}
				res = res.Add(decimal.NewFromFloat(float64(tens * 10)))
				return res, nil
			}
		}
	}

	// 1 - 20
	if idx, err := arrayIndex(s, countNamesSvSE); err == nil {
		return decimal.NewFromFloat(float64(idx)), nil
	}

	if len(s) > 2 {
		last2 := s[len(s)-2:]
		if last2 == ":a" || last2 == ":e" {
			firstPart := s[:len(s)-2]
			return ParseNumber(firstPart)
		}
	}

	return res, fmt.Errorf("error")
}

func parseCountEnglish(s string) (decimal.Decimal, error) {
	if idx, err := arrayIndex(s, countNamesEnUS); err == nil {
		return decimal.NewFromFloat(float64(idx)), nil
	}

	if len(s) > 3 {
		last3 := s[len(s)-3:]
		if last3 == ":th" || last3 == ":rd" {
			firstPart := s[:len(s)-3]
			num, err := decimal.NewFromString(firstPart)
			if err == nil {
				return num, nil
			}
		}
	}

	return decimal.NewFromFloat(0), fmt.Errorf("error")
}
