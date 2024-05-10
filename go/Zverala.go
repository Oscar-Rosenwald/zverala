// -*- eval: (hs-minor-mode 1); -*-
package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

func calculate_a(doubleyear int) int {
	return doubleyear%9 + 1
}

func calculate_b(doubleyearDigits []int, isOutward bool) int {
	add := func(num, thing int) int {
		num += thing
		if num < 0 {
			num = 9
		}
		if num > 9 {
			num = 1
		}
		return num
	}

	var upAngleSlice []int
	for i, d := range doubleyearDigits {
		coefficient := 1
		if i%2 == 0 {
			coefficient = -1
		}
		upAngleSlice = append(upAngleSlice, add(d, coefficient))
	}

	var downAngleSlice []int
	for i, d := range doubleyearDigits {
		coefficient := -1
		if i%2 == 0 {
			coefficient = 1
		}
		downAngleSlice = append(downAngleSlice, add(d, coefficient))
	}

	upAngle := getYearNumberFromSlice(upAngleSlice)
	downAngle := getYearNumberFromSlice(downAngleSlice)
	var partialAnswer float64

	if isOutward {
		partialAnswer = float64(upAngle) / float64(downAngle)
	} else {
		partialAnswer = float64(downAngle) / float64(upAngle)
	}

	precision := math.Pow(10, float64(len(doubleyearDigits)))
	fmt.Printf("partial answer before rounding: %v, precision %v\n", partialAnswer, precision)
	partialAnswer = math.Floor(partialAnswer * precision)

	return addDigits(getYearDigits(int(partialAnswer)))
}

func calculate_c(doubleyear int, doubleyearDigits []int) int {
	digits := make([]*numWithIndex, len(doubleyearDigits))
	for i, d := range doubleyearDigits {
		entry := &numWithIndex{num: d, index: i + 1}
		digits[i] = entry
	}

	// Sort digits by size
	for i := range digits {
		for j := i + 1; j < len(digits); j++ {
			if digits[j].smaller(digits[i]) {
				swapNumbers(digits[i], digits[j])
			}
		}
	}

	bigIndex := getYearNumber(digits)
	coefficient := bigIndex * doubleyear

	return coefficient%9 + 1
}

func getSins(a int) []float64 {
	ret := make([]float64, NUM_CREATURES)

	// sin(ax^2) = 0
	// ax^2 = k*PI
	// x = sqrt( (k*PI) / a )
	for k := 0; k < NUM_CREATURES; k++ {
		root := (float64(k) * float64(math.Pi)) / float64(a)
		res := math.Sqrt(root)
		ret[k] = res
	}

	return ret
}

// TODO Change the floats here and in getSins() to big.Float

func getCosins(b, c int) []float64 {
	// cos(b * sin(cx)) = 0
	// b * sin(cx) = (2k+1 * PI) / 2
	// sin(cx) = (2k+1 * PI) / 2b
	//			 => A
	// cx = sin^-1(A)     -- sin has 2 results for this value!!
	// x = sin^-1(A) / c	   AND
	//     (A - sin^-1(A)) / c
	//
	// HOWEVER!!!!
	//
	// This is only applicable til A < 1, because sin^-1 fails after that.
	// We need to know how many k's there can be, only do those, and then
	// add the cosine period to each of the results.
	//
	// Then there is the period. The period of
	//     sin(a*x)
	// is
	//     2*PI / a.
	// The period of the cosine is already covered, so we add 2*PI / a to
	// all the results given by the various k's.
	//
	// SIMPLIFICATION:
	//
	// Cosine is symmetric about the y-axis, which means we need only consider
	// k >= 0 - we only care about positive results of X. For all the k < 0,
	// it is the case that the results overlap with those of k >= 1.
	//
	// In other works, the resuts repeat for negative k's, which means we can
	// forget about them.
	//
	// Because of this, we can also leave out the 2 in the 2*PI / a period.
	// We'll only get half the results, but the other half will be equal to
	// the first, except for the negative results, which again we don't care for.

	var roots []float64
	for k := 0; true; k++ {
		root := float64(2*k+1) / float64(2*b)
		root = root * float64(math.Pi)
		if root > 1 {
			break
		}
		roots = append(roots, root)
	}

	var baseResults []float64
	for k := 0; k < len(roots); k++ {
		root := roots[k]
		arcsin1 := math.Asin(root)
		arcsin2 := math.Pi - math.Asin(root)

		res1 := arcsin1 / float64(c)
		res2 := arcsin2 / float64(c)

		baseResults = append(baseResults, res1)
		baseResults = append(baseResults, res2)
	}

	sort.Float64s(baseResults)

	var ret []float64
	period := float64(math.Pi) / float64(c)

	for i := 0; i < NUM_CREATURES; i++ {
		fmt.Printf("Animal number: %d\n", i)
		for j := 0; j < len(baseResults); j++ {
			base := baseResults[j]
			newResult := base + period*float64(i)
			fmt.Printf("new result: %v; base result: %v\n", newResult, base)
			ret = append(ret, newResult)
		}
		fmt.Println()
	}
	fmt.Println("==========")
	fmt.Println()

	sort.Float64s(ret)
	return unduplicateSlice(ret)
}

func getOrderedSteps(a, b, c int, sins, cosins []float64) []float64 {
	results := []float64{}
	prev := float64(0)
	sinIndex := 1 // First intersection with x-axis is always from the sine half of the results
	cosIndex := 0
	fmt.Printf("final difference 0 (index 0, previous N/A, smaller 0, s N/A, c N/A, sin index 0, cos index N/A)\n")

	// Start from index 1, because index 0 is always 0. The last one is left
	// out, because that is Chimera, whose days are calculated differently.
	for i := 1; i < NUM_CREATURES; i++ {
		s, c := math.MaxFloat64, math.MaxFloat64
		if len(sins) > sinIndex {
			s = sins[sinIndex]
		}
		if len(cosins) > cosIndex {
			c = cosins[cosIndex]
		}

		var smaller float64
		if s < c {
			smaller = s
			sinIndex++
		} else {
			smaller = c
			cosIndex++
		}

		newResult := smaller - prev
		results = append(results, newResult)
		prev = smaller
	}

	return results
}

func stepsToDays(steps []float64, yearLength, totalSteps float64) []int {
	var days []int
	quiotient := yearLength / totalSteps

	for _, step := range steps {
		ds := int(quiotient * step)
		days = append(days, ds)
	}

	return days
}

func getCreaturesInOrder(direction dirType, chimeraDays int, days []int) []Creature {
	var creatures []Creature

	switch direction {
	case OUT:
		days = append([]int{chimeraDays}, days...)
		for i := 0; i < NUM_CREATURES; i++ {
			b := Creatures[i]
			creatures = append(creatures, Creature{
				name: b,
				days: days[i],
			})
		}
	case IN:
		days = append(days, chimeraDays)
		for i := NUM_CREATURES - 1; i >= 0; i-- {
			creatures = append(creatures, Creature{
				name: Creatures[i],
				days: days[i],
			})
		}
	}

	return creatures
}

func addDragonDays(dragonYear bool, direction dirType, creatures []Creature) []Creature {
	if !dragonYear {
		return creatures
	}

	var retCreatures []Creature
	dragonsAfterIndex := DragonsAfterCreatureIndex
	dragons := Dragons
	lastCreatureIndex := 0

	if direction == IN {
		dragonsAfterIndex = []int{}
		dragons = reverseSlice(Dragons)
		reversedIndexes := reverseSlice(DragonsAfterCreatureIndex)

		for _, index := range reversedIndexes {
			dragonsAfterIndex = append(dragonsAfterIndex, NUM_CREATURES-index-2)
		}
		fmt.Printf("Inversed indexes: %v; dragons: %v\n", dragonsAfterIndex, dragons)
	}

	for i, creatureIndex := range dragonsAfterIndex {
		tailIndex := creatureIndex + 1
		creaturesBefore := creatures[lastCreatureIndex:tailIndex]
		lastCreatureIndex = tailIndex
		newCreatures := make([]Creature, len(creaturesBefore))
		copy(newCreatures, creaturesBefore)
		newCreatures = append(newCreatures, Creature(dragons[i]))
		retCreatures = append(retCreatures, newCreatures...)
	}

	retCreatures = append(retCreatures, creatures[lastCreatureIndex:]...)

	return retCreatures
}

func main() {
	parseArgs()
	_, lastSolstice, nextSolstice := requestYearInfo()
	doubleYear := computeDoubleyear(lastSolstice, nextSolstice)

	if detail, yearKnown := readYearFromFile(doubleYear); yearKnown {
		fmt.Printf("%s\n", detail)
		os.Exit(0)
	}

	dragonYear := isDragonYear(doubleYear)
	doubleYear.dragonYear = dragonYear
	if dragonYear {
		// If this is a dragon year, then we must not consider the days that are
		// reserved for the dragons.
		doubleYear.length -= NUM_DRAGONS
	}

	a := calculate_a(doubleYear.doubleyear)
	b := calculate_b(doubleYear.doubleyearDigits, doubleYear.direction == OUT)
	c := calculate_c(doubleYear.doubleyear, doubleYear.doubleyearDigits)

	sins := getSins(a)
	cosins := getCosins(b, c)

	var (
		orderedSteps          = getOrderedSteps(a, b, c, sins, cosins)
		totalSteps            = sumFloatSlice(orderedSteps)
		days                  = stepsToDays(orderedSteps, float64(doubleYear.length), totalSteps)
		daysSum               = sumIntSlice(days)
		longestCreatureName   = findMax(Creatures, Creatures[0], func(c CreatureName) int { return len(c) })
		maxCreatureNameLength = len(longestCreatureName)
		maxDays               = findMax(days, 0, func(x int) int { return x })
		chimeraDays           = doubleYear.length - daysSum
		padToColumn           = maxCreatureNameLength + 3
		maxDaysLength         = daysDigits(maxDays) + 1
	)

	creaturesInOrder := getCreaturesInOrder(doubleYear.direction, chimeraDays, days)
	creaturesInOrder = addDragonDays(doubleYear.dragonYear, doubleYear.direction, creaturesInOrder)

	printCreatures(creaturesInOrder, doubleYear, padToColumn, maxDaysLength)
	// TODO Add optional logging of the parameters as they are computed.
	// TODO Print full Klvanistic date next to the normal date

	// Good years to test dragons: 2048 (for OUTWARD) and 2049 (for INWARD).
}
