package natural

import (
	"fmt"
	"log"
	"strings"
)

// PresentSvSE returns textual presentation in swedish of input number (5 = "fem")
func PresentSvSE(n int64) string {
	if n == 0 {
		return ""
	}
	s, err := presentSV(n)
	if err != nil {
		log.Fatal(err)
	}
	s = strings.TrimSpace(s) // HACK to remove trailing space from presentSV()
	return s
}

// PresentEnUS returns textual presentation in english of input number (5 = "five")
func PresentEnUS(n int64) string {
	if n == 0 {
		return ""
	}
	s, err := presentEN(n)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func presentSV(n int64) (string, error) {

	// 0 - 19
	if n < 20 {
		for val, i := range numbersToTwentySvSE {
			if i == n {
				return val, nil
			}
		}
		return "", fmt.Errorf("FAILED to find %d", n)
	}

	// 20 - 99
	if n < 100 {
		tiotal := n / 10
		ental := n % 10
		return tensSvSE[tiotal] + PresentSvSE(ental), nil
	}

	// 100 - 999
	if n < 1000 {
		hundratal := n / 100
		last2 := n % 100
		return PresentSvSE(hundratal) + "hundra" + PresentSvSE(last2), nil
	}

	// 1,000 - 999,999
	if n < 1000000 {
		hi := n / 1000
		last3 := n % 1000
		if hi == 1 {
			return "ettusen" + PresentSvSE(last3), nil
		}
		pad := ""
		if hi < 10 {
			pad = " "
		}
		return PresentSvSE(hi) + pad + "tusen " + PresentSvSE(last3), nil
	}

	// 1,000,000 -> 999,999,999 (miljoner)
	if n < 1000000000 {
		millions := n / 1000000
		rest := n % 1000000
		if millions == 1 {
			return "enmiljon" + PresentSvSE(rest), nil
		}
		return PresentSvSE(millions) + " miljoner " + PresentSvSE(rest), nil
	}

	// 1,000,000,000 -> 999,999,999,999 (miljarder)
	if n < 1000000000000 {
		millions := n / 1000000000
		rest := n % 1000000000
		if millions == 1 {
			return "enmiljard" + PresentSvSE(rest), nil
		}
		return PresentSvSE(millions) + " miljarder " + PresentSvSE(rest), nil
	}

	return "", fmt.Errorf("presentSV: FIXME handle %d", n)
}

func presentEN(num int64) (string, error) {

	if num < 20 {
		for val, i := range numbersToTwentyEnUS {
			if i == num {
				return val, nil
			}
		}
		return "", fmt.Errorf("FAILED to find %d", num)
	}

	// 20 - 99
	if num < 100 {
		tiotal := num / 10
		ental := num % 10
		res := tensEnUS[tiotal]
		if ental > 0 {
			res += "-" + PresentEnUS(ental)
		}
		return res, nil
	}

	// 100 - 999
	if num < 1000 {
		hundratal := num / 100
		last2 := num % 100
		res := PresentEnUS(hundratal) + " hundred"
		if last2 > 0 {
			sub, err := presentEN(last2)
			if err != nil {
				return res, err
			}
			res += " and " + sub
		}
		return res, nil
	}

	// 1,000 - 999,999
	if num < 1000000 {
		hi := num / 1000
		last3 := num % 1000
		res := PresentEnUS(hi) + " thousand"
		if last3 > 0 {
			sub, err := presentEN(last3)
			if err != nil {
				return res, err
			}
			res += " " + sub
		}
		return res, nil
	}

	return "", fmt.Errorf("too big")
}

// PresentListSvSE presents a list of strings as "a, b och c"
func PresentListSvSE(list []string) string {
	last := list[len(list)-1]
	rest := list[0 : len(list)-1]
	return strings.Join(rest, ", ") + " och " + last
}

// PresentListEnUS presents a list of strings as "a, b and c"
func PresentListEnUS(list []string) string {
	last := list[len(list)-1]
	rest := list[0 : len(list)-1]
	return strings.Join(rest, ", ") + " and " + last
}
