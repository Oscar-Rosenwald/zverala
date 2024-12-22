package spiral

import (
	"fmt"
	"math"
	"sort"
	ktime "zverala/klvanistic_time"
	"zverala/utils"
)

type (
	kYear      = ktime.KYear
	doubleYear = ktime.DoubleYear
	Creature   = utils.Creature
)

func Calculate_a(doubleyear int) int {
	return doubleyear%9 + 1
}

func Calculate_b(doubleyearDigits []int, isOutward bool) int {
	// add adds NUM and THING. But must be below 10. If the result is > 9,
	// return 1. If the result is < 0, return 9. Otherwise return result.
	add := func(num, thing int) int {
		num += thing
		if num <= 0 {
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

	upAngle := ktime.GetYearNumberFromSlice(upAngleSlice)
	downAngle := ktime.GetYearNumberFromSlice(downAngleSlice)
	utils.PrintDebug("Úhel odrazu: %d", upAngle)
	utils.PrintDebug("Úhel dopadu: %d", downAngle)
	var partialAnswer float64

	if isOutward {
		partialAnswer = float64(upAngle) / float64(downAngle)
	} else {
		partialAnswer = float64(downAngle) / float64(upAngle)
	}

	precision := math.Pow(10, float64(len(doubleyearDigits)))
	partialAnswer = math.Floor(partialAnswer * precision)

	return utils.AddDigits(utils.GetYearDigits(int(partialAnswer)))
}

func Calculate_c(doubleyear int, doubleyearDigits []int) int {
	// 1. Index all digits of the doubleyear number.
	// 2. Order digits in ascending order.
	// 3. Read the Big Index: The number resulting from reading the indexes of the ordered digits.
	// 4. Add the doubleyear number to the big index.
	// 5. c = sum above % 9 + 1

	digits := make([]*utils.NumWithIndex, len(doubleyearDigits))
	for i, d := range doubleyearDigits {
		entry := &utils.NumWithIndex{Num: d, Index: i + 1}
		digits[i] = entry
	}

	// Sort digits by size
	for i := range digits {
		for j := i + 1; j < len(digits); j++ {
			if digits[j].Smaller(digits[i]) {
				utils.SwapNumbers(digits[i], digits[j])
			}
		}
	}

	var digitNums []utils.NumWithIndex
	for _, dig := range digits {
		digitNums = append(digitNums, *dig)
	}
	utils.PrintDEBUG("Cifry kroku podle velikosti: %+v", digitNums)
	bigIndex := ktime.GetYearNumber(digits)
	utils.PrintDebug("Velký index: %d", bigIndex)
	coefficient := bigIndex * doubleyear

	return coefficient%9 + 1
}

func getSins(a int) []float64 {
	utils.PrintDEBUG("Kořeny sinů:")
	ret := make([]float64, utils.NUM_CREATURES)

	// sin(ax^2) = 0
	// ax^2 = k*PI
	// x = sqrt( (k*PI) / a )
	for k := 0; k < utils.NUM_CREATURES; k++ {
		root := (float64(k) * float64(math.Pi)) / float64(a)
		res := math.Sqrt(root)
		ret[k] = res
		utils.PrintDEBUG("SIN: %v", res)
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
	// This is only applicable till A < 1, because sin^-1 fails after that. We
	// need to know how many k's there can be, only do those, and then add the
	// cosine period to each of the results.
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

	utils.PrintDEBUG("Kořeny cosinů:")
	const COS = "COS: "

	var roots []float64
	for k := 0; true; k++ {
		root := float64(2*k+1) / float64(2*b)
		root = root * float64(math.Pi)
		if root > 1 {
			break
		}
		utils.PrintDEBUG(COS+"Základní kořen cosinu: %v", root)
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
	if utils.DebugDebug {
		for _, res := range baseResults {
			utils.PrintDEBUG(COS+"Základní kořen vnitřního sinu: %v", res)
		}
	}

	var ret []float64
	period := float64(math.Pi) / float64(c)

	for i := 0; i < utils.NUM_CREATURES; i++ {
		for j := 0; j < len(baseResults); j++ {
			base := baseResults[j]
			newResult := base + period*float64(i)
			ret = append(ret, newResult)
		}
	}

	sort.Float64s(ret)
	if utils.DebugDebug {
		for _, res := range ret {
			utils.PrintDEBUG(COS+"Nový kořen: %v", res)
		}
	}

	return utils.UnduplicateSlice(ret)
}

// GetOrderedSteps composes a list of ordered floats where each float
// corresponds to a length of Lysak's jump. Sins and cosins must both be of
// length utils.NUM_CREATURES and be ordered. The result has length utils.NUM_CREATURES - 1
// (where the 1 is reserved for Chimera).
func GetOrderedSteps(sins, cosins []float64) []float64 {
	utils.PrintDebug("Počítám déklů Lysákových skoků")
	results := []float64{}
	prev := float64(0)
	sinIndex := 1 // First intersection with x-axis is always from the sine half of the results
	cosIndex := 0

	printStep := func(length float64, index int, sin, cos, smaller float64, sinIndex, cosIndex int) {
		sinStr := fmt.Sprintf("%v", sin)
		if sin < 0 {
			sinStr = "N/A"
		}
		cosStr := fmt.Sprintf("%v", cos)
		if cos < 0 {
			cosStr = "N/A"
		}

		utils.PrintDEBUG("Dékla skoku: %v (index bysosti: %d; sin: %s; cos: %s; sin index: %d; cos index: %d; menší: %v)",
			length, index, sinStr, cosStr, sinIndex, cosIndex, smaller)
	}

	printStep(0, 0, -1, -1, 0, 0, 0)

	// Start from index 1, because index 0 is always 0. The last one is left
	// out, because that is Chimera, whose days are calculated differently.
	for i := 1; i < utils.NUM_CREATURES; i++ {
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
		printStep(newResult, i, s, c, smaller, sinIndex, cosIndex)
	}

	return results
}

// stepsToDays converts lengths of Lysak's jumps to number of days along
// Zverala's spiral.
func stepsToDays(steps []float64, yearLength, totalSteps float64) []int {
	var days []int
	quiotient := yearLength / totalSteps

	for _, step := range steps {
		ds := int(quiotient * step)
		days = append(days, ds)
		utils.PrintDEBUG("Lysákův skok dlouhý %v koresponduje počtem dnů %d", step, ds)
	}

	return days
}

// GetCreaturesInOrder gives a list of creatures and their corresponding day
// ranges during the year. The order of creatures is given by the direction of
// the kyear. utils.Dragons are ignored.
func GetCreaturesInOrder(direction ktime.DirType, chimeraDays int, days []int) []utils.Creature {
	var creatures []utils.Creature

	switch direction {
	case ktime.OUT:
		days = append([]int{chimeraDays}, days...)
		for i := 0; i < utils.NUM_CREATURES; i++ {
			b := utils.Creatures[i]
			creatures = append(creatures, utils.Creature{
				Name: b,
				Days: days[i],
			})
		}
	case ktime.IN:
		days = append(days, chimeraDays)
		for i := utils.NUM_CREATURES - 1; i >= 0; i-- {
			creatures = append(creatures, utils.Creature{
				Name: utils.Creatures[i],
				Days: days[len(days)-i-1],
			})
		}
	}

	utils.PrintDEBUG("Bytosti v pořadí (bez draků): %v", creatures)

	return creatures
}

// AddDragonDays takes the output of getCreaturesInOrder and adds dragons to it
// if dragonYear is true.
func AddDragonDays(dragonYear bool, direction ktime.DirType, creatures []utils.Creature) []utils.Creature {
	if !dragonYear {
		return creatures
	}

	utils.PrintDebug("Přidávám draky k bytostem")
	var retCreatures []utils.Creature
	dragonsAfterIndex := utils.DragonsAfterCreatureIndex
	dragons := utils.Dragons
	lastCreatureIndex := 0

	if direction == ktime.IN {
		dragonsAfterIndex = []int{}
		dragons = utils.ReverseSlice(utils.Dragons)
		reversedIndexes := utils.ReverseSlice(utils.DragonsAfterCreatureIndex)

		for _, index := range reversedIndexes {
			dragonsAfterIndex = append(dragonsAfterIndex, utils.NUM_CREATURES-index-2)
		}
	}

	utils.PrintDEBUG("Draci následují po bytostech s těmito indexy: %v", dragonsAfterIndex)

	for i, creatureIndex := range dragonsAfterIndex {
		tailIndex := creatureIndex + 1
		creaturesBefore := creatures[lastCreatureIndex:tailIndex]
		lastCreatureIndex = tailIndex
		newCreatures := make([]utils.Creature, len(creaturesBefore))
		copy(newCreatures, creaturesBefore)
		newCreatures = append(newCreatures, utils.Creature(dragons[i]))
		retCreatures = append(retCreatures, newCreatures...)
	}

	retCreatures = append(retCreatures, creatures[lastCreatureIndex:]...)
	utils.PrintDebug("Bytosti s draky: %v", retCreatures)

	return retCreatures
}

// ComputeZverala returns a list of Creatures in the order they appear in kyear
// along with how long they are active. The daysPerCreature return value is for
// convenience of the caller; they needn't be in order.
func ComputeZverala(kyear *kYear, doubleyear *doubleYear) (creaturesInOrder []Creature, daysPerCreature []int) {
	dragonYear := ktime.IsDragonYear(*kyear)
	kyear.DragonYear = dragonYear
	if dragonYear {
		utils.PrintDebug("%d je krok draků", kyear.NormalYearStart.Year())
		// If this is a dragon year, then we must not consider the days that are
		// reserved for the dragons.
		kyear.Length -= utils.NUM_DRAGONS
	}

	utils.PrintDebug("Propočítávám pohybové zákony")
	a := Calculate_a(kyear.Doubleyear)
	utils.PrintDebug("Parametr a: %d", a)
	b := Calculate_b(kyear.DoubleyearDigits, kyear.Direction == ktime.OUT)
	utils.PrintDebug("Parametr b: %d", b)
	c := Calculate_c(kyear.Doubleyear, kyear.DoubleyearDigits)
	utils.PrintDebug("Parametr c: %d", c)

	sins := getSins(a)
	cosins := getCosins(b, c)

	orderedSteps := GetOrderedSteps(sins, cosins)
	totalSteps := utils.SumFloatSlice(orderedSteps)
	daysPerCreature = stepsToDays(orderedSteps, float64(kyear.Length), totalSteps)
	daysSum := utils.SumIntSlice(daysPerCreature)
	chimeraDays := kyear.Length - daysSum

	creaturesInOrder = GetCreaturesInOrder(kyear.Direction, chimeraDays, daysPerCreature)
	creaturesInOrder = AddDragonDays(kyear.DragonYear, kyear.Direction, creaturesInOrder)

	return creaturesInOrder, append(daysPerCreature, chimeraDays)
}
