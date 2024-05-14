// -*- eval: (hs-minor-mode 1); -*-
package main

import (
	"fmt"
	"math"
	"time"
)

type char = byte
type dirType string

const (
	OUT dirType = "OUTWARD"
	IN  dirType = "INWARD"
)

func (typ dirType) toChar() char {
	switch typ {
	case OUT:
		return 'O'
	case IN:
		return 'S'
	}
	// Pass an error to hanldeError so we guarantee that the program is
	// terminated here.
	handleError(fmt.Errorf("Invalid K. year direction: '%s'", typ))
	return ' '
}

type kYear struct {
	doubleyear       int
	doubleyearDigits []int
	direction        dirType
	normalYearStart  time.Time
	dragonYear       bool
	length           int
}

func (k kYear) toString() string {
	return fmt.Sprintf("%s - %d %s (dragons:%t) (%v) (length %d)\n", k.normalYearStart.String(), k.doubleyear, k.direction, k.dragonYear, k.doubleyearDigits, k.length)
}

func (k kYear) toReadableString() string {
	return fmt.Sprintf("Rok %d / %d %c - %d days", k.normalYearStart.Year()+1, k.doubleyear, k.direction.toChar(), k.length)
}

func getYearDigits(doubleyear int) []int {
	digit := doubleyear % 10
	if doubleyear < 10 {
		return []int{digit}
	}
	return append(getYearDigits(doubleyear/10), digit)
}

var saveToFile = true
var file = "../Zverala2.txt" // TODO Change for the appropriate file.
var referenceTime kYear = kYear{
	doubleyear:       27073,
	doubleyearDigits: []int{2, 7, 0, 7, 3},
	normalYearStart:  time.Date(2013, 12, 21, 11, 0, 0, 0, time.UTC), // 12:00 in Prague
	direction:        OUT,
	dragonYear:       false,
}

type doubleYear struct {
	outKyear kYear
	inKyear  kYear
	endTime  time.Time
	length   int
}

func getYearNumber(digits []*numWithIndex) int {
	var ret = 0
	var reverse []numWithIndex

	for _, d := range digits {
		reverse = append([]numWithIndex{*d}, reverse...)
	}

	for i, d := range reverse {
		ret += int(math.Pow(10, float64(i))) * d.index
	}

	return ret
}

func getYearNumberFromSlice(digits []int) int {
	var ret = 0
	var reverse []int

	for _, d := range digits {
		reverse = append([]int{d}, reverse...)
	}

	for i, d := range reverse {
		ret += int(math.Pow(10, float64(i))) * d
	}

	return ret
}

// computeKyear returns the current kyear denoted by yearStart and yearEnd.
func computeKyear(yearStart, yearEnd time.Time) kYear {
	yearDiff := int(math.Abs(float64(yearStart.Year() - referenceTime.normalYearStart.Year())))
	beforeReference := yearStart.Before(referenceTime.normalYearStart)
	evenDiff := yearDiff%2 == 0

	var dir = IN
	if evenDiff {
		dir = OUT
	}

	var kyearDiff = yearDiff / 2
	if beforeReference {
		kyearDiff = (yearDiff + 1) / -2
	}

	doubleyear := referenceTime.doubleyear + kyearDiff

	yearLength := int(math.Abs(float64(yearStart.Sub(yearEnd) / (time.Hour * 24))))

	return kYear{
		doubleyear:       doubleyear,
		doubleyearDigits: getYearDigits(doubleyear),
		direction:        dir,
		normalYearStart:  yearStart,
		length:           yearLength,
	}
}

// isDragonYear reports whether the doubleyear of which kyear is part contains
// dragons.
func isDragonYear(kyear kYear) bool {
	printDebug("Je to krok draga?")
	dragonDigit := kyear.doubleyearDigits[len(kyear.doubleyearDigits)-1]
	var (
		dragonNum_3 = 3
		dragonNum_5 = 5
		dragonNum_7 = 7
	)

	if kyear.direction == OUT {
		dragonNum_5 = dragonDigit
	} else {
		dragonNum_7 = dragonDigit
	}

	printDebug("Dračí čísla jsou: %d %d %d", dragonNum_3, dragonNum_5, dragonNum_7)

	dy := kyear.doubleyear
	divisibleBy := 0

	doForHolyNumber := func(num int) {
		if num != 0 && dy%num == 0 {
			divisibleBy++
		}
	}

	doForHolyNumber(dragonNum_3)
	doForHolyNumber(dragonNum_5)
	doForHolyNumber(dragonNum_7)

	return divisibleBy == 2
}

type kDate struct {
	doubleYear doubleYear
	dir        dirType
	// 1 to 7 plus a dummy 8th for the dragon days.
	sun int
	// 0 for the 0 day of the sun, 1 to 17 for the remaining vyks. If sun == 8,
	// ignored.
	vyk kVyk
	// - 0 if vyk == 0 and sun != 8. This is the 0 day of each sun.
	// - 1 to 9 (ish) during dragon days (when sun == 8).
	day kDay
}

type kDay int

const (
	KDAY_STREDA kDay = 1
	KDAY_TESLEK kDay = 2
	KDAY_TRETEK kDay = 3
	KDAY_DELE   kDay = 4
	KDAY_VYKON  kDay = 5
	KDAY_VYKOFF kDay = 6
)

func (d kDay) toString() string {
	switch d {
	case 0:
		return "nultý den"
	case KDAY_STREDA:
		return "středa"
	case KDAY_TESLEK:
		return "teslek"
	case KDAY_TRETEK:
		return "třetek"
	case KDAY_DELE:
		return "děle"
	case KDAY_VYKON:
		return "výkon"
	case KDAY_VYKOFF:
		return "výkoff"
	default:
		return "neznámý den"
	}

}

type kVyk int

func (vyk kVyk) toOrderedNumber() string {
	switch vyk {
	case 1:
		return "prvního"
	case 2:
		return "druhého"
	case 3:
		return "třetího"
	case 4:
		return "čtvrtého"
	case 5:
		return "pátého"
	case 6:
		return "šestého"
	case 7:
		return "sedmého"
	case 8:
		return "osmého"
	case 9:
		return "devátého"
	case 10:
		return "desátého"
	case 11:
		return "jedenáctého"
	case 12:
		return "dvanáctého"
	case 13:
		return "třináctého"
	case 14:
		return "čtrnáctého"
	case 15:
		return "patnáctého"
	case 16:
		return "šestnáctého"
	case 17:
		return "sedmnáctého"
	default:
		return "neznámého"
	}
}

type Planet string

var Planets = []Planet{
	"Pluta (chaosu)",
	"Neptunu (chaosu)",
	"Uranu (vzduchu)",
	"Saturnu (smrti)",
	"Jupitera (dřeva)",
	"Marsu (vody)",
	"Země (života)",
	"Venuše (země)",
	"Merkuru (ohně)",
}

func (kd *kDate) toString() string {
	kYear := kd.doubleYear.outKyear
	if kd.dir == IN {
		kYear = kd.doubleYear.inKyear
	}

	if kd.sun > NUM_SUNS {
		// TODO The document claims that there can only be 7-9 dargon days, but
		// counting reveals there could be more. Did we do research on this, or
		// is it a mistake?

		nonDragonDays := kd.doubleYear.length - NUM_SUNS*SUN_LENGTH
		var startPlanetIndex int
		switch nonDragonDays {
		case 7:
			startPlanetIndex = 2
		case 8:
			startPlanetIndex = 1
		case 9:
			startPlanetIndex = 0
		default:
			handleError(fmt.Errorf("Invalidní počet dní planet: %d", nonDragonDays))
		}

		planet := Planets[startPlanetIndex+int(kd.day)-1]
		return fmt.Sprintf("krok %d %c, %-34s",
			kYear.doubleyear,
			kYear.direction.toChar(),
			fmt.Sprintf("den %s", string(planet)),
		)
	}

	if kd.vyk == 0 {
		// Zero day
		return fmt.Sprintf("krok %d %c, %d. slunce, %-23s",
			kYear.doubleyear,
			kYear.direction.toChar(),
			kd.sun,
			"nultý den",
		)
	}

	return fmt.Sprintf("krok %d %c, %d. slunce, %-23s",
		kYear.doubleyear,
		kYear.direction.toChar(),
		kd.sun,
		fmt.Sprintf("%s %s výku", kd.day.toString(), kd.vyk.toOrderedNumber()),
	)
}

// A sun has a 0 day and then 17 vyks with 6 days each. All suns are the same
// length.
const SUN_LENGTH = 103
const VYK_LENGTH = 6
const NUM_SUNS = 7

// timeToKlvanisticDate transforms a specific time into a klvanistic date.
//
// doubleYearStart and doubleYearEnd are the timestamps of winter solstices at
// the ends of the doubleyear. Note that they must be two years apart - they
// denote a doubleyear, not a kyear. MidPoint is the noon on the winter solstice
// in the middle of the doubleyear.
//
// Date must occur between doubleYearStart and doubleYearEnd. It must point to
// midnight on that day.
func timeToKlvanisticDate(date time.Time, doubleYear doubleYear) kDate {
	midPoint := doubleYear.inKyear.normalYearStart
	daysFromStart := int(date.Sub(doubleYear.outKyear.normalYearStart).Hours() / 24)

	printDEBUG("Hledám datum %d ve dvojroku %d, který má %d dní. Datum je %d dní od začátku.",
		date.Format("2.1. 2006"),
		doubleYear.inKyear.doubleyear,
		doubleYear.length,
		daysFromStart)

	suns := daysFromStart/SUN_LENGTH + 1
	daysInSun := daysFromStart%SUN_LENGTH + 1

	// Always subtract 1 from daysInSun, because all suns start with a zero day.
	vyks := (daysInSun - 1) / VYK_LENGTH
	daysInVyk := (daysInSun - 1) % VYK_LENGTH
	if daysInVyk == 0 {
		daysInVyk = VYK_LENGTH
	}
	printDEBUG("Sluncí: %d, dní ve slunci: %d, výků: %d, dní ve výku: %d", suns, daysInSun, vyks, daysInVyk)

	if suns == NUM_SUNS+1 {
		// It's the end of the doubleyear. This isn't a real sun, it's the
		// dragon days.

		printDEBUG("Datum je den draka")
		vyks = 0              // No vyks during dragon days.
		daysInVyk = daysInSun // Ignore the zero day during dragon days.
	}

	var dir dirType
	if date.After(midPoint) {
		dir = IN
	} else {
		dir = OUT
	}

	return kDate{
		doubleYear: doubleYear,
		dir:        dir,
		sun:        suns,
		vyk:        kVyk(vyks),
		day:        kDay(daysInVyk),
	}
}
