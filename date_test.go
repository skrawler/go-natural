package natural

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMonthDay(t *testing.T) {
	md := NewMonthDay("12-15")
	assert.Equal(t, "12-15", md.String())
}

func TestMonthDayMarshal(t *testing.T) {
	md := NewMonthDay("12-15")
	d, err := md.MarshalJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, "12-15", string(d))
}

func TestMonthDayUnmarshal(t *testing.T) {
	d := []byte("12-15")
	md := MonthDay{}
	err := md.UnmarshalJSON(d)
	assert.Equal(t, nil, err)
	assert.Equal(t, "12-15", md.String())
}

func TestParseDateIntoMonthDay(t *testing.T) {
	md, err := ParseDateIntoMonthDay("15 december")
	assert.Equal(t, nil, err)
	assert.Equal(t, "12-15", md.String())
}

func TestParseYear(t *testing.T) {
	y, err := ParseYear("")
	assert.NotEqual(t, nil, err)

	expected := map[string]int{
		"1900": 1900,
		"98":   1998,
		"45":   1945,
		"12":   2012,
	}
	for s, i := range expected {
		y, err = ParseYear(s)
		assert.Equal(t, nil, err)
		assert.Equal(t, i, y)
	}
}

func TestParseMonth(t *testing.T) {
	m, err := ParseMonth("")
	assert.NotEqual(t, nil, err)

	expected := map[string]time.Month{
		// eng
		"jan":  time.January,
		"may":  time.May,
		"june": time.June,
		// swe
		"januari": time.January,
		"oktober": time.October,
	}
	for s, i := range expected {
		m, err = ParseMonth(s)
		assert.Equal(t, nil, err)
		assert.Equal(t, i, m)
	}
}

func TestParseWeekday(t *testing.T) {
	m, err := ParseWeekday("")
	assert.NotEqual(t, nil, err)

	expected := map[string]time.Weekday{
		// eng
		"saturday":  time.Saturday,
		"Wednesday": time.Wednesday,
		"Sat":       time.Saturday,
		// swe
		"Lördag":   time.Saturday,
		"lördagen": time.Saturday,
		"lör":      time.Saturday,
	}
	for s, i := range expected {
		m, err = ParseWeekday(s)
		assert.Equal(t, nil, err)
		assert.Equal(t, i, m)
	}
}

func TestRelative(t *testing.T) {
	_, err := ParseTime("")
	assert.NotEqual(t, nil, err)

	_, err = ParseTime("the 28:th of ma")
	assert.NotEqual(t, nil, err)

	// eng
	t1, err := ParseTime("the 28:th of march")
	assert.Equal(t, nil, err)
	assert.Equal(t, "the 28:th of March", readableDateNoYear(t1, "en_US"))

	// swe
	t1, err = ParseTime("den 28:e mars")
	assert.Equal(t, nil, err)
	assert.Equal(t, "the 28:th of March", readableDateNoYear(t1, "en_US"))

	expected := map[string]string{
		// swe
		"6":                       "06:00",
		"sex":                     "06:00",
		"middag":                  "12:00",
		"midnatt":                 "00:00",
		"sex på morgonen":         "06:00",
		"sex på kvällen":          "18:00",
		"kvart i sju":             "06:45",
		"kvart över sju":          "07:15",
		"halv elva":               "10:30",
		"halv elva på kvällen":    "22:30",
		"tjugo över elva":         "11:20",
		"tjugo minuter över elva": "11:20",
		"tjugo i elva":            "10:40",
		"tjugo minuter i elva":    "10:40",
		"arton och trettio":       "18:30",
		"kl 18":                   "18:00",
		"klockan 18:30":           "18:30",
		"18":                      "18:00",
		"18:33":                   "18:33",
		"18:33:59":                "18:33",
		"idag":                    "00:00",
		"igår":                    "INNAN 00:00", // XXX igår
		"imorgon":                 "INNAN 00:00", // XXX i morgon
	}

	for s, i := range expected {
		t1, err = ParseTime(s)
		assert.Equal(t, nil, err)
		assert.Equal(t, i, relativeReadableTime(t1))
	}

	/*

	   $this->assertEquals(Carbon::parse("first saturday of march")->toDateString(), CarbonSwedish::parse("den första lördagen i mars")->toDateString());
	   $this->assertEquals(Carbon::parse("third sunday of march")->toDateString(), CarbonSwedish::parse("den tredje söndagen i mars")->toDateString());
	   $this->assertEquals(Carbon::parse("first saturday of march 2014")->toDateString(), CarbonSwedish::parse("den första lördagen i mars 2014")->toDateString());

	   $this->assertEquals(Carbon::parse("first saturday of march 2015")->toDateString(), CarbonSwedish::parse("1:a lördag i mars 2015")->toDateString());
	   $this->assertEquals(Carbon::parse("first saturday of march 2015")->toDateString(), CarbonSwedish::parse("1:a lör i mar 2015")->toDateString());

	   $this->assertEquals(Carbon::parse("2015-03-31")->toDateString(), CarbonSwedish::parse("sista dagen i mars 2015")->toDateString());


	   // FIXME
	   //$this->assertEquals(Carbon::parse("2015-03-31 16:20:00")->toDateTimeString(), CarbonSwedish::parse("sista dagen i mars 2015 kl 16:20")->toDateSTimetring());
	   //$this->assertEquals(Carbon::parse("2015-03-31 16:00:00")->toDateTimeString(), CarbonSwedish::parse("sista dagen i mars kl 16")->toDateTimeString());

	   $this->assertEquals(Carbon::parse("this sunday")->toDateString(), CarbonSwedish::parse("på söndag")->toDateString());

	   $this->assertEquals(Carbon::parse("sunday next week")->toDateString(), CarbonSwedish::parse("nästa söndag")->toDateString());


	   $this->assertEquals(Carbon::parse("june 2008")->toDateString(), CarbonSwedish::parse("juni 2008")->toDateString());
	   // PHP BUG ?: "june" is not set to day 1 in php 5.4, 5.5 or 5.6
	   $this->assertEquals(Carbon::parse("june")->toDateString(), CarbonSwedish::parse("juni")->toDateString());

	   $this->assertEquals(Carbon::parse("first day of january 2008")->toDateString(), CarbonSwedish::parse("första dagen i januari 2008")->toDateString());
	   $this->assertEquals(Carbon::parse("first day of january")->toDateString(), CarbonSwedish::parse("första dagen i januari")->toDateString());

	   // PHP BUG?: "fifth day of march 2015" dont seem to parse :(
	   $this->assertEquals(Carbon::parse("2015-03-05")->toDateString(), CarbonSwedish::parse("femte dagen i mars 2015")->toDateString());

	   $this->assertEquals(Carbon::parse("last day of march")->toDateString(), CarbonSwedish::parse("sista dagen i mars")->toDateString());


	   // PHP BUG? "second month 2008" = 2015-07-06 ????
	   $this->assertEquals(Carbon::parse("2008-02-01")->toDateString(), CarbonSwedish::parse("andra månaden 2008")->toDateString());
	   $this->assertEquals(Carbon::parse(date("Y")."-03-01")->toDateString(), CarbonSwedish::parse("tredje månaden")->toDateString());
	*/
}

/*
func TestExistingFunctionality()
{
    // make sure english strings pass as before
    $this->assertEquals("2015-05-09", CarbonSwedish::parse("may 9, 2015")->toDateString());

    $this->assertEquals("2012-02-05", CarbonSwedish::parse("2012-02-05")->toDateString());
}

// TODO: förra tisdagen
// TODO: kommande tisdag = nästa tisdag
// TODO måndag nästa vecka
*/

/*
public function testParseWithYear()
{
    $this->assertEquals("2015-01-01", CarbonSwedish::parse("1 januari 2015")->toDateString());
    $this->assertEquals("2015-05-02", CarbonSwedish::parse("2 maj, 2015")->toDateString());
    $this->assertEquals("2015-05-03", CarbonSwedish::parse("3 maj,2015")->toDateString());
    $this->assertEquals("2015-04-01", CarbonSwedish::parse("1 apr, 2015")->toDateString());
    $this->assertEquals("2015-04-02", CarbonSwedish::parse("2:a apr, 2015")->toDateString());
    $this->assertEquals("2015-05-05", CarbonSwedish::parse("5:e maj, 2015")->toDateString());
    $this->assertEquals("2015-05-05", CarbonSwedish::parse("femte maj, 2015")->toDateString());
}

public function testParseWithoutYear()
{
    $this->assertEquals(date("Y")."-01-01", CarbonSwedish::parse("1 januari")->toDateString());
    $this->assertEquals(date("Y")."-05-29", CarbonSwedish::parse("29 maj")->toDateString());
    $this->assertEquals(date("Y")."-04-29", CarbonSwedish::parse("29 apr")->toDateString());
    $this->assertEquals(date("Y")."-09-02", CarbonSwedish::parse("2:a sep")->toDateString());
    $this->assertEquals(date("Y")."-10-05", CarbonSwedish::parse("5:e okt")->toDateString());
    $this->assertEquals(date("Y")."-05-05", CarbonSwedish::parse("femte maj")->toDateString());
    $this->assertEquals(date("Y")."-05-05 18:24:00", CarbonSwedish::parse("femte maj 18:24:00")->toDateTimeString());
    $this->assertEquals(date("Y")."-05-05", CarbonSwedish::parse("den femte maj")->toDateString());
    $this->assertEquals(date("Y")."-02-03", CarbonSwedish::parse("den 3:e feb")->toDateString());
}

public function testCompound()
{
    // https://php.net/manual/en/datetime.formats.compound.php

    $this->assertEquals("2015-05-05 19:31:10", CarbonSwedish::parse("femte maj 19:31:10, 2015")->toDateTimeString());

    $this->assertEquals("2008-02-01 14:30:00", CarbonSwedish::parse("den 1:a feb 14:30 2008")->toDateTimeString());

    $this->assertEquals(Carbon::parse(date("Y")."-02-01 14:30:00")->toDateTimeString(), CarbonSwedish::parse("den 1:a feb 14:30")->toDateTimeString());

    $this->assertEquals(Carbon::parse(date("Y")."-02-01 14:30:00")->toDateTimeString(), CarbonSwedish::parse("den 1:a feb kl 14:30")->toDateTimeString());

    $this->assertEquals(Carbon::parse(date("Y")."-02-01 14:00:00")->toDateTimeString(), CarbonSwedish::parse("den 1:a feb kl 14")->toDateTimeString());


    $this->assertEquals(Carbon::parse("this sunday 16:00")->toDateTimeString(), CarbonSwedish::parse("på söndag kl 16:00")->toDateTimeString());
    $this->assertEquals(Carbon::parse("this sunday 16:00")->toDateTimeString(), CarbonSwedish::parse("på söndag 16:00")->toDateTimeString());
    $this->assertEquals(Carbon::parse("this sunday 16:00")->toDateTimeString(), CarbonSwedish::parse("på söndag kl 16")->toDateTimeString());


    $this->assertEquals(Carbon::parse("tomorrow 18:00")->toDateTimeString(), CarbonSwedish::parse("imorgon 18:00")->toDateTimeString());
    $this->assertEquals(Carbon::parse("yesterday 18:00")->toDateTimeString(), CarbonSwedish::parse("yesterday 18:00")->toDateTimeString());


    $this->assertEquals(Carbon::parse("tomorrow noon")->toDateTimeString(), CarbonSwedish::parse("imorgon middag")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow midnight")->toDateTimeString(), CarbonSwedish::parse("imorgon natt")->toDateTimeString());


    $this->assertEquals(Carbon::parse("tomorrow 18:00")->toDateTimeString(), CarbonSwedish::parse("imorgon klockan sex på kvällen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 06:00")->toDateTimeString(), CarbonSwedish::parse("imorgon klockan sex på morgonen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 06:00")->toDateTimeString(), CarbonSwedish::parse("imorgon klockan sex")->toDateTimeString());

    $this->assertEquals(Carbon::parse("tomorrow 18:30")->toDateTimeString(), CarbonSwedish::parse("imorgon kl arton och trettio")->toDateTimeString());

    $this->assertEquals(Carbon::parse("tomorrow 10:30")->toDateTimeString(), CarbonSwedish::parse("imorgon halv elva")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 10:30")->toDateTimeString(), CarbonSwedish::parse("imorgon halv elva på morgonen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 22:30")->toDateTimeString(), CarbonSwedish::parse("imorgon halv elva på kvällen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 22:45")->toDateTimeString(), CarbonSwedish::parse("imorgon kvart i elva på kvällen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 11:15")->toDateTimeString(), CarbonSwedish::parse("imorgon kvart över elva på morgonen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 23:15")->toDateTimeString(), CarbonSwedish::parse("imorgon kvart över elva på kvällen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 23:20")->toDateTimeString(), CarbonSwedish::parse("imorgon tjugo över elva på kvällen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 23:08")->toDateTimeString(), CarbonSwedish::parse("imorgon åtta över elva på kvällen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 18:55")->toDateTimeString(), CarbonSwedish::parse("imorgon fem i sju på kvällen")->toDateTimeString());
    $this->assertEquals(Carbon::parse("tomorrow 18:55")->toDateTimeString(), CarbonSwedish::parse("imorgon fem minuter i sju på kvällen")->toDateTimeString());
}

public function testTimeOnly()
{
    $this->assertEquals(date("Y-m-d")." 19:31:10", CarbonSwedish::parse("19:31:10")->toDateTimeString());
    $this->assertEquals(date("Y-m-d")." 19:30:00", CarbonSwedish::parse("19:30")->toDateTimeString());
    $this->assertEquals(date("Y-m-d")." 19:30:00", CarbonSwedish::parse("kl 19:30")->toDateTimeString());
    $this->assertEquals(date("Y-m-d")." 19:30:00", CarbonSwedish::parse("klockan 19:30")->toDateTimeString());

    $this->assertEquals(date("Y-m-d")." 21:00:00", CarbonSwedish::parse("kl 21")->toDateTimeString());
}

*/

// used for testing
func readableDateNoYear(t time.Time, locale string) string {

	switch locale {
	case "sv_SE":
		return fmt.Sprintf("den %d:e %s", t.Day(), t.Month())
	}

	return fmt.Sprintf("the %d:th of %s", t.Day(), t.Month())
}

// used for testing
func relativeReadableTime(t time.Time) string {

	// XXX if today
	now := time.Now()
	if dateEquals(now, t) {
		return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
	}

	if dateLessThan(now, t) {
		return fmt.Sprintf("INNAN %02d:%02d", t.Hour(), t.Minute())
	}
	return fmt.Sprintf("EFTER %02d:%02d", t.Hour(), t.Minute())
}

func dateEquals(t1 time.Time, t2 time.Time) bool {
	if t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day() {
		return true
	}
	return false
}

// dateLessThan returns true if t1 is before t2 in the calendar
func dateLessThan(t1 time.Time, t2 time.Time) bool {
	d1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	d2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())
	return d1.Unix() < d2.Unix()
}
