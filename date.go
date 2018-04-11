package natural

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseWeekday parses a weekday name into a time.Weekday
func ParseWeekday(s string) (time.Weekday, error) {
	s = ucFirst(s)
	if val, ok := WeekdayNames[s]; ok {
		return val, nil
	}
	if len(s) > 2 && s[len(s)-2:] == "en" {
		// swe: lördagen -> lördag
		return ParseWeekday(s[:len(s)-2])
	}
	return 0, fmt.Errorf("Cannot parse weekday: %s", s)
}

// ParseTime parses a string like HH:MM, HH:MM:SS, "klockan sex på kvällen" etc into a time.Time
func ParseTime(s string) (time.Time, error) {
	t := time.Now()
	if s == "" {
		return t, fmt.Errorf("empty")
	}

	t = setMinute(t, 0)
	t = setSecond(t, 0)

	if num, err := strconv.ParseInt(s, 10, 64); err == nil {
		return setHour(t, num), nil
	}

	if s == "middag" || s == "lunch" {
		return setHour(t, 12), nil
	}

	if s == "midnatt" || s == "natt" {
		return setHour(t, 0), nil
	}

	timeBase := int64(0)

	// https://sv.wikipedia.org/wiki/F%C3%B6rmiddag
	re := regexp.MustCompile(`^(?P<time>[\w\d\s]+)+ (?:på morgonen|på förmiddagen|förmiddag|fm)+$`)
	match := re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		s = match[0][1]
	}

	// https://sv.wikipedia.org/wiki/Eftermiddag
	re = regexp.MustCompile(`^(?P<time>[\w\d\s]+)+ (på kvällen|i kväll|på eftermiddagen|eftermiddag|em)+$`)
	match = re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		timeBase = 12
		s = match[0][1]
	}

	if s == "idag" || s == "i dag" {
		return beginningOfDay(t), nil
	}

	if s == "igår" || s == "i går" {
		return subDay(t, 1), nil
	}

	if s == "imorgon" || s == "i morgon" || s == "imorrn" || s == "imorron" || s == "i morron" {
		return addDay(t, 1), nil
	}

	// "18:23:59", "18:23", "18", "kl 18:30"
	re = regexp.MustCompile(`^(?:kl |klockan )?(?P<hour>[\d]+)+:?(?P<min>[\d]+)*:?(?P<sec>[\d:]+)*$`)
	match = re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		_hr, err := ParseNumber(match[0][1])
		if err != nil {
			return t, err
		}
		hr := _hr.IntPart()
		t = setHour(t, timeBase+hr)
		if match[0][2] != "" {
			_mn, err := ParseNumber(match[0][2])
			mn := _mn.IntPart()
			if err != nil {
				return t, err
			}
			t = setMinute(t, mn)
		}
		if match[0][3] != "" {
			_sc, err := ParseNumber(match[0][3])
			sc := _sc.IntPart()
			if err != nil {
				return t, err
			}
			t = setSecond(t, sc)
		}
		return t, nil
	}

	re = regexp.MustCompile(`^kvart i (?P<time>[\w\d]+)+$`)
	match = re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		_hr, err := ParseNumber(match[0][1])
		hr := _hr.IntPart()
		if err == nil {
			t = setHour(t, timeBase+hr-1)
			t = setMinute(t, 45)
			return t, nil
		}
	}

	re = regexp.MustCompile(`^kvart över (?P<time>[\w\d]+)+$`)
	match = re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		_hr, err := ParseNumber(match[0][1])
		hr := _hr.IntPart()
		if err == nil {
			t = setHour(t, timeBase+hr)
			t = setMinute(t, 15)
			return t, nil
		}
	}

	// "halv elva", "halv elva på morgonen"
	re = regexp.MustCompile(`^halv (?P<time>[\w\d\s]+)+$`)
	match = re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		_hr, err := ParseNumber(match[0][1])
		hr := _hr.IntPart()
		if err == nil {
			t = setHour(t, timeBase+hr-1)
			t = setMinute(t, 30)
			return t, nil
		}
	}

	// "tjugo över elva", "tjugo minuter över elva"
	re = regexp.MustCompile(`^(?P<min>[\w\d]+)+(?: minuter| min)? över (?P<time>[\w\d]+)+$`)
	match = re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		_mn, err := ParseNumber(match[0][1])
		mn := _mn.IntPart()
		if err == nil {
			t = setMinute(t, mn)
			_hr, err := ParseNumber(match[0][2])
			hr := _hr.IntPart()
			if err == nil {
				t = setHour(t, timeBase+hr)
				return t, nil
			}
		}
	}

	// "tjugo i elva", "tjugo minuter i elva"
	re = regexp.MustCompile(`^(?P<min>[\w\d]+)+ (?:minuter |min )?i (?P<time>[\w\d]+)+$`)
	match = re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		_mm, err := ParseNumber(match[0][1])
		mm := _mm.IntPart()
		if err == nil {
			t = setMinute(t, 60-mm)
			_hr, err := ParseNumber(match[0][2])
			hr := _hr.IntPart()
			if err == nil {
				t = setHour(t, timeBase+hr-1)
				return t, nil
			}
		}
	}

	// "arton och trettio"
	re = regexp.MustCompile(`^(?P<hour>[\w\d]+)+ och (?P<min>[\w\d]+)+$`)
	match = re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		_mm, err := ParseNumber(match[0][2])
		mm := _mm.IntPart()
		if err == nil {
			t = setMinute(t, mm)
			_hr, err := ParseNumber(match[0][1])
			hr := _hr.IntPart()
			if err == nil {
				if hr > 12 {
					timeBase = 0
				}
				t = setHour(t, timeBase+hr)
				return t, nil
			}
		}
	}

	// "den 28:e mars", "the 28:th of may"
	re = regexp.MustCompile(`^(?:den |the )?(?P<day>[\w:]+)? (?:of )?(?P<month>[\w]+)\s?(kl |klockan |at )?(?P<time>[:\d]+)?(?:,?)\s?(?P<year>[0-9]+)?$`)
	match = re.FindAllStringSubmatch(s, -1)
	if len(match) != 0 {
		dd, err := ParseCount(match[0][1])
		if err != nil {
			return t, err
		}
		t = setDay(t, dd.IntPart())
		month, err := ParseMonth(match[0][2])
		if err != nil {
			return t, err
		}
		t = setMonth(t, month)

		// XXX TODO parse time of day
		return t, nil
	}

	// ex "sex"
	_hr, err := ParseNumber(s)
	if err == nil {
		hr := _hr.IntPart()
		t = setHour(t, timeBase+hr)
		return t, nil
	}

	return t, fmt.Errorf("failed to parse: %s", s)
}

/*
func parseRelativeTimeSvSE(when string) *time.Time {

    t := time.Now()

    tz := nil

    if when == "nu" {
        return &t
    }



    re := regexp.MustCompile("^(imorgon|i morgon|imorrn|imorron|i morron)+ (kl |klockan )?" + timeRegex + "$")

    match := re.FindAllStringSubmatch(when, -1)
    d(match)


    preg_match("/^(imorgon|i morgon|imorrn|imorron|i morron)+ (kl |klockan )?".$timeRegex."$/ui", $time, $match);
    if (!empty($match["time"])) {
        $res->addDay(1);
        return self::parseTimeToObject($match["time"], $res);
    }
*/

/*
    preg_match("/^(igår|i går)+ (kl |klockan )?".$timeRegex."$/ui", $time, $match);
    if (!empty($match["time"])) {
        $res->subDay(1);
        return self::parseTimeToObject($match["time"], $res);
    }

    // "den andra månaden"
    preg_match("/^(den )?(?<count>[\:\w]+)? månaden(,?)\s?(?<year>[0-9]+)?$/ui", $time, $match);
    if (!empty($match["count"])) {
        $count = self::parseCount($match["count"]);

        if (!empty($match["year"])) {
            $res->year = self::parseYear($match["year"]);
        }

        if ($match["count"] == "sista") {
            $res->month = 12;
        } else {
            $res->month = $count;
        }
        $res->day = 1;
        return $res;
    }

    // "den första lördagen i mars", "den sista dagen i mars"
    preg_match("/^(den )?(?<count>[\:\w]+)? (?<weekday>\w+) i (?<month>\w+)(,?)\s?(?<year>[0-9]+)?$/ui", $time, $match);
    if (!empty($match["count"]) && !empty($match["weekday"]) && !empty($match["month"])) {
//d($match);
        if (!empty($match["year"])) {
            $res->year = self::parseYear($match["year"]);
        }

        $res->month = self::parseMonth($match["month"]);

        if ($match["count"] == "sista" && $match["weekday"] == "dagen") {

            $res->day = $res->daysInMonth;
        } else {
            $res->day = self::parseCount($match["count"]);
        }

        if ($match["weekday"] == "dagen") {
            return $res;
        }

        $weekday = self::parseWeekday($match["weekday"]);
        if ($weekday === null) {
            return $res;
        }
//                d($match);
        $foundCount = 0;
        for ($i = 1; $i <= $res->daysInMonth; $i++) {
            $res->day = $i;
            if ($res->dayOfWeek == $weekday) {
                $foundCount++;
                if ($foundCount == self::parseCount($match["count"])) {
                    return $res;
                }
            }
        }
    }


    // "6 maj, 2015", "6:e maj 2015", "6 maj 2015", "6 maj", "6:e maj", "sjätte maj", "6 maj 14:30:00, 2015"
    preg_match("/^(den )?(?<day>[\w:]+)? (?<month>[\w]+)\s?(kl |klockan )?".$timeRegex."(,?)\s?(?<year>[0-9]+)?$/ui", $time, $match);

    if (!empty($match["day"]) && !empty($match["month"])) {
        $day = self::parseCount($match["day"]);
        if ($day) {
            $res->day = $day;
            $month = self::parseMonth($match["month"]);
            if ($month) {
                $res->month = $month;

                if (!empty($match["time"]) && is_numeric($match["time"])) {
                    $match["time"] = (int) $match["time"];
                    if ($match["time"] >= 1000 &&  $match["time"] < 2900) {
                        $match["year"] = $match["time"];
                        unset($match["time"]);
                    }
                }

//                    d($match);
                if (!empty($match["year"])) {
                    $res->year = self::parseYear($match["year"]);
                }
                if (!empty($match["time"])) {
                    $res = self::parseTimeToObject($match["time"], $res);
                }
                return $res;
            }
        }
    }

    // "juni", "juni 2008"
    preg_match("/^(?<month>\w+)(,?)\s?(?<year>[0-9]+)?$/ui", $time, $match);
    if (!empty($match["month"])) {
        $month = self::parseMonth($match["month"]);
        $res->month = $month;
        if (!empty($match["year"])) {
            $res->year = self::parseYear($match["year"]);
            // NOTE: if year is also specified, reset day to 1 to mimic DateTime behaviour
            $res->day = 1;
        }
        return $res;
    }

    // "på onsdag", "på onsdag kl 16:00" (denna onsdag)
    preg_match("/^på (?<weekday>\w+)\s?(kl |klockan )?".$timeRegex."$/ui", $time, $match);
    if (!empty($match["weekday"])) {
        $weekday = self::parseWeekday($match["weekday"]);
        if ($weekday !== null) {
            for ($i = 0; $i <= 7; $i++) {
                $res->addDays(1);
                if ($res->dayOfWeek == $weekday) {
                    if (!empty($match["time"])) {
                        return self::parseTimeToObject($match["time"], $res);
                    }
                    return $res;
                }
            }
        }
    }

    // "nästa onsdag", "nästa torsdag kl 16:00" (närmaste torsdagen, efter idag)
    preg_match("/^nästa (?<weekday>\w+)\s?(kl |klockan )?".$timeRegex."$/ui", $time, $match);
    if (!empty($match["weekday"])) {
        $weekday = self::parseWeekday($match["weekday"]);
        $res->addDays(7);
        if ($weekday !== null) {
            for ($i = 0; $i <= 7; $i++) {
                $res->addDays(1);
                if ($res->dayOfWeek == $weekday) {
                    if (!empty($match["time"])) {
                        return self::parseTimeToObject($match["time"], $res);
                    }
                    return $res;
                }
            }
        }
    }

    return parent::parse($time, $tz);
}
*/

// ParseMonth turns textual representation into a time.Month
func ParseMonth(s string) (time.Month, error) {

	s = ucFirst(s)

	// swe
	if month, ok := MonthNames[s]; ok {
		return month, nil
	}

	// eng
	for m := 1; m <= 12; m++ {
		if time.Month(m).String() == s {
			return time.Month(m), nil
		}
	}

	return 0, fmt.Errorf("Cannot parse month: %s", s)
}

// ParseYear parses a 2 or 4 digit year string into a int
func ParseYear(s string) (int, error) {

	year, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	if year < 100 {
		// get last 2 digits of current year
		t := time.Now()
		currentYear, _, _ := t.Date()
		yearStr := fmt.Sprintf("%d", currentYear)
		max, err := strconv.ParseInt(yearStr[2:], 10, 64)
		if err != nil {
			return 0, err
		}
		if year < max {
			return int(2000 + year), nil
		}
		return int(1900 + year), nil
	}

	return int(year), nil
}

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func addDay(t time.Time, diff int) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day+diff, 0, 0, 0, 0, t.Location())
}

func subDay(t time.Time, diff int) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day+diff, 0, 0, 0, 0, t.Location())
}

func setMonth(t time.Time, month time.Month) time.Time {
	year, _, day := t.Date()
	return time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

func setDay(t time.Time, day int64) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, int(day), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

func setHour(t time.Time, hour int64) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, int(hour), t.Minute(), t.Second(), 0, t.Location())
}

func setMinute(t time.Time, min int64) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, t.Hour(), int(min), t.Second(), 0, t.Location())
}

func setSecond(t time.Time, sec int64) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, t.Hour(), t.Minute(), int(sec), 0, t.Location())
}

// MonthDay ...
type MonthDay struct {
	Month time.Month
	Day   int64
}

// formats December 15 into "12-15" (MM-DD) format
func (md MonthDay) String() string {
	return fmt.Sprintf("%02d-%02d", md.Month, md.Day)
}

// MarshalJSON ...
func (md MonthDay) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"` + md.String() + `"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON ...
func (md *MonthDay) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	*md = NewMonthDay(s)
	return nil
}

// NaturalSwe returns "13:e December"
func (md *MonthDay) NaturalSwe() string {
	return PresentCountShortSwedish(int(md.Day)) + " " + MonthsSvSE[md.Month]
}

// NewMonthDay ...
func NewMonthDay(s string) MonthDay {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		log.Fatal("invalid input " + s)
	}
	day, err := ParseNumber(parts[1])
	if err != nil {
		log.Fatal(err)
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatal(err)
	}
	md := MonthDay{}
	md.Day = day.IntPart()
	md.Month = time.Month(month)
	if err != nil {
		log.Fatal(err)
	}
	return md
}

// ParseDateIntoMonthDay parses "15 december" into "12-15" (MM-DD) format
func ParseDateIntoMonthDay(s string) (MonthDay, error) {
	md := MonthDay{}
	var err error

	parts := strings.SplitN(s, " ", 2)
	if len(parts) != 2 {
		return md, fmt.Errorf("not 2 parts")
	}

	md.Month, err = ParseMonth(parts[1])
	if err != nil {
		return md, err
	}

	day, err := ParseNumber(parts[0])
	if err != nil {
		return md, err
	}

	md.Day = day.IntPart()
	return md, nil
}
