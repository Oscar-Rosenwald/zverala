// -*- eval: (hs-minor-mode 1); -*-
package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

func handleError(err error) {
	if err != nil {
		_, f, no, _ := runtime.Caller(1)
		fmt.Printf("Eror '%s' na řádku %s:%d\n", err, f, no)
		os.Exit(1)
	}
}

func addDigits(digits []int) int {
	if len(digits) == 1 {
		return digits[0]
	}

	var sum = 0
	for _, d := range digits {
		sum += d
	}
	sumDigits := getYearDigits(sum)
	return addDigits(sumDigits)
}

type numWithIndex struct {
	num   int
	index int
}

func (nwi *numWithIndex) smaller(other *numWithIndex) bool {
	return nwi.num != other.num && nwi.num < other.num
}

func unPointer(digits []*numWithIndex) []numWithIndex {
	var ret []numWithIndex
	for _, d := range digits {
		ret = append(ret, *d)
	}
	return ret
}

func swapNumbers(one, two *numWithIndex) {
	helpNum := one.num
	helpIndex := one.index

	one.num = two.num
	one.index = two.index

	two.num = helpNum
	two.index = helpIndex
}

// Assumes the slice is sorted
func unduplicateSlice(slice []float64) []float64 {
	if len(slice) == 0 {
		return nil
	}

	var ret []float64
	ret = append(ret, slice[0])

	// Start from 1 as we always will include the first item
	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[i-1] {
			ret = append(ret, slice[i])
		}
	}

	return ret
}

func sumIntSlice(s []int) int {
	var ret int
	for _, el := range s {
		ret += el
	}
	return ret
}

func sumFloatSlice(s []float64) float64 {
	var ret float64
	for _, el := range s {
		ret += el
	}
	return ret
}

func findMax[T any](s []T, def T, getSize func(T) int) T {
	for _, el := range s {
		size := getSize(el)
		if size > getSize(def) {
			def = el
		}
	}
	return def
}

func pad(maxLength int, l int, padChar string) string {
	return strings.Repeat(padChar, maxLength-l)
}

// We must consider the apparent length of the string, so we count the
// number of runes, not the number of bytes.
func strLen(s string) int {
	return utf8.RuneCountInString(s)
}

func daysDigits(i int) int {
	return len(strconv.Itoa(i))
}

func dayString(days int) string {
	switch days {
	case 0:
		return "dnů"
	case 1:
		return "den"
	case 2, 3, 4:
		return "dny"
	default:
		return "dní"
	}
}

func reverseSlice[T any](s []T) []T {
	var ret []T
	for i := 0; i < len(s); i++ {
		ret = append(ret, s[len(s)-i-1])
	}
	return ret
}
