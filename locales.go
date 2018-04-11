package natural

import "time"

// ...
var (
	tensSvSE = []string{
		"---", "tio", "tjugo", "trettio", "fyrtio", "femtio",
		"sextio", "sjuttio", "åttio", "nittio",
	}

	tensEnUS = []string{
		"---", "ten", "twenty", "thirty", "forty", "fifty",
		"sixty", "seventy", "eighty", "ninety",
	}

	numbersToTwentySvSE = map[string]int64{
		"ett":     1,
		"två":     2,
		"tre":     3,
		"fyra":    4,
		"fem":     5,
		"sex":     6,
		"sju":     7,
		"åtta":    8,
		"nio":     9,
		"tio":     10,
		"elva":    11,
		"tolv":    12,
		"tretton": 13,
		"fjorton": 14,
		"femton":  15,
		"sexton":  16,
		"sjutton": 17,
		"arton":   18,
		"nitton":  19,
		"tjugo":   20,
	}

	numbersToTwentyEnUS = map[string]int64{
		"one":       1,
		"two":       2,
		"three":     3,
		"four":      4,
		"five":      5,
		"six":       6,
		"seven":     7,
		"eight":     8,
		"nine":      9,
		"ten":       10,
		"eleven":    11,
		"twelve":    12,
		"thirteen":  13,
		"fourteen":  14,
		"fifteen":   15,
		"sixteen":   16,
		"seventeen": 17,
		"eighteen":  18,
		"nineteen":  19,
		"twenty":    20,
	}

	countNamesSvSE = []string{
		"---",
		"första", "andra", "tredje", "fjärde", "femte",
		"sjätte", "sjunde", "åttonde", "nionde", "tionde",
		"elfte", "tolfte", "trettonde", "fjortonde", "femtonde",
		"sextonde", "sjuttonde", "artonde", "nittonde", "tjugonde",
		/*
			"tjugoförsta", "tjugoandra", "tjugotredje", "tjugofjärde",
			"tjugofemte", "tjugosjätte", "tjugosjunde", "tjugoåttonde",
			"tjugonionde", "trettionde", "trettioförsta",
		*/
	}

	countNamesEnUS = []string{
		"---",
		"first", "second", "third", "fourth", "fifth",
		"sixth", "seventh", "eighth", "ninth", "tenth",
		"eleventh", "twelfth", "thirteenth", "fourteenth", "fifteenth",
		"sixteenth", "seventeenth", "eighteenth", "nineteenth", "twentieth",
	}

	WeekdayNames = map[string]time.Weekday{
		// eng
		"Sunday":    time.Sunday,
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
		"Sun":       time.Sunday,
		"Mon":       time.Monday,
		"Tue":       time.Tuesday,
		"Wed":       time.Wednesday,
		"Thu":       time.Thursday,
		"Fri":       time.Friday,
		"Sat":       time.Saturday,

		// swe
		"Söndag":  time.Sunday,
		"Måndag":  time.Monday,
		"Tisdag":  time.Tuesday,
		"Onsdag":  time.Wednesday,
		"Torsdag": time.Thursday,
		"Fredag":  time.Friday,
		"Lördag":  time.Saturday,
		"Sön":     time.Sunday,
		"Mån":     time.Monday,
		"Tis":     time.Tuesday,
		"Ons":     time.Wednesday,
		"Tor":     time.Thursday,
		"Fre":     time.Friday,
		"Lör":     time.Saturday,
	}

	MonthNames = map[string]time.Month{
		// eng
		"January":   time.January,
		"February":  time.February,
		"March":     time.March,
		"April":     time.April,
		"May":       time.May,
		"June":      time.June,
		"July":      time.July,
		"August":    time.August,
		"September": time.September,
		"October":   time.October,
		"November":  time.November,
		"December":  time.December,
		"Jan":       time.January,
		"Feb":       time.February,
		"Mar":       time.March,
		"Apr":       time.April,
		"Jun":       time.June,
		"Jul":       time.July,
		"Aug":       time.August,
		"Sep":       time.September,
		"Oct":       time.October,
		"Nov":       time.November,
		"Dec":       time.December,

		// swe
		"Januari":  time.January,
		"Februari": time.February,
		"Mars":     time.March,
		"Maj":      time.May,
		"Juni":     time.June,
		"Juli":     time.July,
		"Augusti":  time.August,
		"Oktober":  time.October,
		"Okt":      time.October,
	}

	WeekdaysSvSE = map[time.Weekday]string{
		time.Monday:    "Måndag",
		time.Tuesday:   "Tisdag",
		time.Wednesday: "Onsdag",
		time.Thursday:  "Torsdag",
		time.Friday:    "Fredag",
		time.Saturday:  "Lördag",
		time.Sunday:    "Söndag",
	}

	MonthsSvSE = map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Mars",
		time.April:     "April",
		time.May:       "Maj",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Augusti",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "December",
	}

	fractionsSvSE = map[string]string{
		"halv": "1/2", "hälften": "1/2",
		"tredjedel":  "1/3",
		"fjärdedel":  "1/4",
		"femtedel":   "1/5",
		"sjättedel":  "1/6",
		"sjundedel":  "1/7",
		"åttondedel": "1/8", "åttondel": "1/8",
		"niondedel": "1/9", "niondel": "1/9",
		"tiondedel": "1/10", "tiondel": "1/10", "tiodel": "1/10",
		"elftedel":     "1/11",
		"tolftedel":    "1/12",
		"trettondedel": "1/13", "trettondel": "1/13",
		"fjortondedel": "1/14", "fjortondel": "1/14",
		"femtondedel": "1/15", "femtondel": "1/15",
		"sextondedel": "1/16", "sextondel": "1/16",
		"sjuttondedel": "1/17", "sjuttondel": "1/17",
		"artondedel": "1/18", "artondel": "1/18",
		"nittondedel": "1/19", "nittondel": "1/19",
		"tjugondedel": "1/20", "tjugondel": "1/20", "tjugodel": "1/20",
		"trettiondedel": "1/30", "trettiondel": "1/30",
		"fyrtiondedel": "1/40", "fyrtiondel": "1/40",
		"femtiondedel": "1/50", "femtiondel": "1/50",
		"sextiondedel": "1/60", "sextiondel": "1/60",
		"sjuttiondedel": "1/70", "sjuttiondel": "1/70",
		"åttiondedel": "1/80", "åttiondel": "1/80",
		"nittiondedel": "1/90", "nittiondel": "1/90",
		"hundradedel": "1/100", "hundradel": "1/100",
	}
)
