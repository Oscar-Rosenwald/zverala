// -*- eval: (hs-minor-mode 1); -*-
package klvanistic_time

import (
	"fmt"
	"math"
	"time"
	"zverala/utils"
)

type Char = byte
type DirType string

const (
	OUT DirType = "OUTWARD"
	IN  DirType = "INWARD"
)

func (typ DirType) ToChar() Char {
	switch typ {
	case OUT:
		return 'O'
	case IN:
		return 'S'
	}
	// Pass an error to hanldeError so we guarantee that the program is
	// terminated here.
	utils.HandleError(fmt.Errorf("Invalidní směr kroku: '%s'", typ))
	return ' '
}

type KYear struct {
	Doubleyear       int
	DoubleyearDigits []int
	Direction        DirType
	NormalYearStart  time.Time
	DragonYear       bool
	Length           int
}

func (k KYear) ToString() string {
	return fmt.Sprintf("%s - %d %s (dragons:%t) (%v) (length %d)\n", k.NormalYearStart.String(), k.Doubleyear, k.Direction, k.DragonYear, k.DoubleyearDigits, k.Length)
}

func (k KYear) ToReadableString() string {
	return fmt.Sprintf("Rok %d / %d %c - %d days", k.NormalYearStart.Year()+1, k.Doubleyear, k.Direction.ToChar(), k.Length)
}

var saveToFile = true
var referenceTime KYear = KYear{
	Doubleyear:       27073,
	DoubleyearDigits: []int{2, 7, 0, 7, 3},
	NormalYearStart:  time.Date(2013, 12, 21, 11, 0, 0, 0, time.UTC), // 12:00 in Prague
	Direction:        OUT,
	DragonYear:       false,
}

type DoubleYear struct {
	OutKyear KYear
	InKyear  KYear
	EndTime  time.Time
	Length   int
}

// ToCache returns dy in the cache form to be written into the cache file.
func (dy *DoubleYear) ToCache() string {
	return fmt.Sprintf("%d:%d:%d:%d:%d\n",
		dy.OutKyear.Doubleyear,
		dy.OutKyear.NormalYearStart.Year(),
		dy.OutKyear.NormalYearStart.Day(),
		dy.InKyear.NormalYearStart.Day(),
		dy.EndTime.Day(),
	)
}

func (dy *DoubleYear) ToString() string {
	start := dy.OutKyear.NormalYearStart
	end := dy.EndTime
	return fmt.Sprintf("Dvojrok %d.%d. %d - %d.%d. %d (%d dní)", start.Day(), start.Month(), start.Year(), end.Day(), end.Month(), end.Year(), dy.Length)
}

// GetYearNumber turns digits into a single int.
func GetYearNumber(digits []*utils.NumWithIndex) int {
	var ret = 0
	var reverse []utils.NumWithIndex

	for _, d := range digits {
		reverse = append([]utils.NumWithIndex{*d}, reverse...)
	}

	for i, d := range reverse {
		ret += int(math.Pow(10, float64(i))) * d.Index
	}

	return ret
}

// GetYearNumberFromSlice turns digits into a single int.
func GetYearNumberFromSlice(digits []int) int {
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

// ComputeKyear returns the current kyear denoted by yearStart and yearEnd.
func ComputeKyear(yearStart, yearEnd time.Time) KYear {
	yearDiff := int(math.Abs(float64(yearStart.Year() - referenceTime.NormalYearStart.Year())))
	beforeReference := yearStart.Before(referenceTime.NormalYearStart)
	evenDiff := yearDiff%2 == 0

	var dir = IN
	if evenDiff {
		dir = OUT
	}

	var kyearDiff = yearDiff / 2
	if beforeReference {
		kyearDiff = (yearDiff + 1) / -2
	}

	doubleyear := referenceTime.Doubleyear + kyearDiff

	yearLength := int(math.Abs(float64(yearStart.Sub(yearEnd) / (time.Hour * 24))))

	return KYear{
		Doubleyear:       doubleyear,
		DoubleyearDigits: utils.GetYearDigits(doubleyear),
		Direction:        dir,
		NormalYearStart:  yearStart,
		Length:           yearLength,
	}
}

// IsDragonYear reports whether the doubleyear of which kyear is part contains
// dragons.
func IsDragonYear(kyear KYear) bool {
	utils.PrintDebug("Je to krok draka?")
	dragonDigit := kyear.DoubleyearDigits[len(kyear.DoubleyearDigits)-1]
	var (
		dragonNum_3 = 3
		dragonNum_5 = 5
		dragonNum_7 = 7
	)

	if kyear.Direction == OUT {
		dragonNum_5 = dragonDigit
	} else {
		dragonNum_7 = dragonDigit
	}

	utils.PrintDebug("Dračí čísla jsou: %d %d %d", dragonNum_3, dragonNum_5, dragonNum_7)

	dy := kyear.Doubleyear
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

// KDate is a representation of a single day in the Klvanistic calendar.
type KDate struct {
	doubleYear DoubleYear
	dir        DirType
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

func (d kDay) ToString() string {
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

func (vyk kVyk) ToOrderedNumber() string {
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

// There are multiple Plutos because it's the planet that can be assigned to
// several days, depending on how many days there are in the doubleyear.
var Planets = []Planet{
	"Pluta (chaosu)",
	"Pluta (chaosu)",
	"Pluta (chaosu)",
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

func (kd *KDate) ToString() string {
	kYear := kd.doubleYear.OutKyear
	if kd.dir == IN {
		kYear = kd.doubleYear.InKyear
	}

	if kd.sun > NUM_SUNS {
		dragonDays := kd.doubleYear.Length - NUM_SUNS*SUN_LENGTH
		if dragonDays < 7 {
			utils.HandleError(fmt.Errorf("Invalidní počet dní planet: %d", dragonDays))
		}

		startPlanetIndex := 5 - (dragonDays - 7)
		utils.PrintDEBUG("Počáteční index planet je %d pro %d. den planet z %d (z totálního počtu dní %d)",
			startPlanetIndex,
			kd.day,
			dragonDays,
			kd.doubleYear.Length)

		planet := Planets[startPlanetIndex+int(kd.day)-1]
		return fmt.Sprintf("krok %d %c, %-34s",
			kYear.Doubleyear,
			kYear.Direction.ToChar(),
			fmt.Sprintf("den %s", string(planet)),
		)
	}

	if kd.vyk == 0 {
		// Zero day
		return fmt.Sprintf("krok %d %c, %d. slunce, %-23s",
			kYear.Doubleyear,
			kYear.Direction.ToChar(),
			kd.sun,
			"nultý den",
		)
	}

	return fmt.Sprintf("krok %d %c, %d. slunce, %-23s",
		kYear.Doubleyear,
		kYear.Direction.ToChar(),
		kd.sun,
		fmt.Sprintf("%s %s výku", kd.day.ToString(), kd.vyk.ToOrderedNumber()),
	)
}

// A sun has a 0 day and then 17 vyks with 6 days each. All suns are the same
// length.
const SUN_LENGTH = 103
const VYK_LENGTH = 6
const NUM_SUNS = 7

// TimeToKlvanisticDate transforms a specific time into a Klvanistic date.
func TimeToKlvanisticDate(date time.Time, doubleYear DoubleYear) KDate {
	midPoint := doubleYear.InKyear.NormalYearStart
	// Here we do +1 because the distance between 21st December and the same day
	// should be taken as 1. That way were actually counting the index of the
	// day from 1.
	daysFromStart := int(date.Sub(doubleYear.OutKyear.NormalYearStart).Hours()/24) + 1

	utils.PrintDEBUG("Hledám datum %s ve dvojroku %d, který má %d dní. Datum je %d dní od začátku.",
		date.Format("2.1. 2006"),
		doubleYear.InKyear.Doubleyear,
		doubleYear.Length,
		daysFromStart)

	var suns, daysInSun, vyks, daysInVyk int

	// Counting the suns is a bit tricky. If date is on the same day as the
	// doubleyear's start (AKA daysFromStart == 0), then suns = 1. If not, we
	// can't just divide. For any number of days with a remainder of 0, dividing
	// and adding 1 would work, but if the remainder is 0, then we mustn't add
	// the 1.
	//
	// Imagine a SUN which is 2 days long. Here are the options:
	//  - 0 days from start of the doubleyear: SUN == 1
	//  - 1 day from the start:                SUN == 1
	//  - 2 days from the start:               SUN == 1
	//  - 3 days from the start:               SUN == 2
	//  - 4 days from the start:               SUN == 2
	//
	// This isn't a straighforward mathematical operation. The 0 case must be
	// checked separately, and then we can subtract 1 from the days, divide, and
	// add one.
	//
	// AKA: if (days == 0) -> sun = 1
	//      else           -> sun = (days - 1) / SUN_LENGTH  + 1
	//
	// Days in the sun are easier. We can take the remainder after dividing by
	// SUN_LENGTH. However, if the remainder is 0, then days of the sun are
	// equal to the length of the sun.
	if daysFromStart == 0 {
		suns = 1
		daysInSun = 1
	} else {
		suns = (daysFromStart-1)/SUN_LENGTH + 1
		daysInSun = daysFromStart % SUN_LENGTH
		if daysInSun == 0 {
			daysInSun = SUN_LENGTH
		}
	}

	if daysInSun == 1 {
		vyks = 0
		daysInVyk = 0
	} else {
		// Minus 1 for zero day, minus another 1 as explained above.
		vyks = (daysInSun-2)/VYK_LENGTH + 1
		daysInVyk = (daysInSun - 1) % VYK_LENGTH
		if daysInVyk == 0 {
			daysInVyk = VYK_LENGTH
		}
	}

	utils.PrintDEBUG("Sluncí: %d, dní ve slunci: %d, výků: %d, dní ve výku: %d", suns, daysInSun, vyks, daysInVyk)

	if suns == NUM_SUNS+1 {
		// It's the end of the doubleyear. This isn't a real sun, it's the
		// dragon days.

		utils.PrintDEBUG("Datum je den draka")
		vyks = 0              // No vyks during dragon days.
		daysInVyk = daysInSun // Ignore the zero day during dragon days.
	}

	var dir DirType
	if date.Before(midPoint) {
		dir = OUT
	} else {
		dir = IN
	}

	return KDate{
		doubleYear: doubleYear,
		dir:        dir,
		sun:        suns,
		vyk:        kVyk(vyks),
		day:        kDay(daysInVyk),
	}
}
