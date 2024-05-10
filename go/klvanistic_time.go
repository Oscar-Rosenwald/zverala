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
var file = "../Zverala2.txt"
var referenceTime kYear = kYear{
	doubleyear:       27073,
	doubleyearDigits: []int{2, 7, 0, 7, 3},
	normalYearStart:  time.Date(2013, 12, 21, 11, 0, 0, 0, time.UTC), // 12:00 in Prague
	direction:        OUT,
	dragonYear:       false,
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

func computeDoubleyear(normalYear, yearEnd time.Time) kYear {
	yearDiff := int(math.Abs(float64(normalYear.Year() - referenceTime.normalYearStart.Year())))
	beforeReference := normalYear.Before(referenceTime.normalYearStart)
	evenDiff := yearDiff%2 == 0

	var dir = IN
	if evenDiff {
		dir = OUT
	}

	var doubleyearDiff = yearDiff / 2
	if beforeReference {
		doubleyearDiff = (yearDiff + 1) / -2
	}

	doubleyear := referenceTime.doubleyear + doubleyearDiff

	yearLength := int(math.Abs(float64(normalYear.Sub(yearEnd) / (time.Hour * 24))))

	return kYear{
		doubleyear:       doubleyear,
		doubleyearDigits: getYearDigits(doubleyear),
		direction:        dir,
		normalYearStart:  normalYear,
		length:           yearLength,
	}
}

func isDragonYear(doubleYear kYear) bool {
	printDebug("Je to krok draga?")
	dragonDigit := doubleYear.doubleyearDigits[len(doubleYear.doubleyearDigits)-1]
	var (
		dragonNum_3 = 3
		dragonNum_5 = 5
		dragonNum_7 = 7
	)

	if doubleYear.direction == OUT {
		dragonNum_5 = dragonDigit
	} else {
		dragonNum_7 = dragonDigit
	}

	printDebug("Dračí čísla jsou: %d %d %d", dragonNum_3, dragonNum_5, dragonNum_7)

	dy := doubleYear.doubleyear
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
