// -*- eval: (hs-minor-mode 1); -*-
package utils

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

func HandleError(err error) {
	if err != nil {
		_, f, no, _ := runtime.Caller(1)
		fmt.Printf("Eror '%s' na řádku %s:%d\n", err, f, no)
		os.Exit(1)
	}
}

// Print helpful debug logs.
var Debug = false

// Print unhelpful debug logs.
var DebugDebug = false

func PrintDebug(msg string, args ...interface{}) {
	if !Debug {
		return
	}

	if len(args) > 0 {
		fmt.Printf(msg+"\n", args...)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

func PrintDEBUG(msg string, args ...interface{}) {
	if !DebugDebug {
		return
	}

	if len(args) > 0 {
		fmt.Printf(msg+"\n", args...)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

func PrintInfo(msg string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf(msg+"\n", args...)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

func GetYearDigits(doubleyear int) []int {
	digit := doubleyear % 10
	if doubleyear < 10 {
		return []int{digit}
	}
	return append(GetYearDigits(doubleyear/10), digit)
}

func AddDigits(digits []int) int {
	if len(digits) == 1 {
		return digits[0]
	}

	var sum = 0
	for _, d := range digits {
		sum += d
	}
	sumDigits := GetYearDigits(sum)
	return AddDigits(sumDigits)
}

type NumWithIndex struct {
	Num   int
	Index int
}

func (nwi *NumWithIndex) Smaller(other *NumWithIndex) bool {
	return nwi.Num != other.Num && nwi.Num < other.Num
}

func unPointer(digits []*NumWithIndex) []NumWithIndex {
	var ret []NumWithIndex
	for _, d := range digits {
		ret = append(ret, *d)
	}
	return ret
}

func SwapNumbers(one, two *NumWithIndex) {
	helpNum := one.Num
	helpIndex := one.Index

	one.Num = two.Num
	one.Index = two.Index

	two.Num = helpNum
	two.Index = helpIndex
}

// Assumes the slice is sorted
func UnduplicateSlice(slice []float64) []float64 {
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

func SumIntSlice(s []int) int {
	var ret int
	for _, el := range s {
		ret += el
	}
	return ret
}

func SumFloatSlice(s []float64) float64 {
	var ret float64
	for _, el := range s {
		ret += el
	}
	return ret
}

func FindMax[T any](s []T, def T, getSize func(T) int) T {
	for _, el := range s {
		size := getSize(el)
		if size > getSize(def) {
			def = el
		}
	}
	return def
}

func Pad(maxLength int, l int, padChar string) string {
	return strings.Repeat(padChar, maxLength-l)
}

// We must consider the apparent length of the string, so we count the
// number of runes, not the number of bytes.
func StrLen(s string) int {
	return utf8.RuneCountInString(s)
}

func DaysDigits(i int) int {
	return len(strconv.Itoa(i))
}

func DayString(days int) string {
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

func ReverseSlice[T any](s []T) []T {
	var ret []T
	for i := 0; i < len(s); i++ {
		ret = append(ret, s[len(s)-i-1])
	}
	return ret
}
