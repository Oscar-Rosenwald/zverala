// -*- eval: (hs-minor-mode 1); -*-
package main

import (
	"zverala/command_line"
	ktime "zverala/klvanistic_time"
	"zverala/spiral"
	"zverala/utils"
)

type doubleYear = ktime.DoubleYear
type kYear = ktime.KYear

func main() {
	command_line.ParseArgs()
	doubleYear, kyear, yearCached := command_line.RequestYearInfo()
	utils.PrintDebug("Zpracovávám krok %s a dvojrok %s", kyear.ToReadableString(), doubleYear.ToString())
	utils.PrintDEBUG("Zpracovávám krok %s a dvojrok %s", kyear.ToString(), doubleYear.ToString())

	creaturesInOrder, daysPerCreature := spiral.ComputeZverala(&kyear, &doubleYear)

	var (
		longestCreatureName   = utils.FindMax(utils.Creatures, utils.Creatures[0], func(c utils.CreatureName) int { return len(c) })
		maxCreatureNameLength = len(longestCreatureName)
		maxDays               = utils.FindMax(daysPerCreature, 0, func(x int) int { return x })
		padToColumn           = maxCreatureNameLength + 3
		maxDaysLength         = utils.DaysDigits(maxDays) + 1
	)

	utils.PrintDebug("Nejdelší jméno bytosti: %d (%s)", maxCreatureNameLength, longestCreatureName)
	utils.PrintDebug("Nejdleší počet dní (v cifrách): %d", maxDaysLength)

	command_line.PrintCreatures(creaturesInOrder, doubleYear, kyear, padToColumn, maxDaysLength)

	if !yearCached {
		command_line.WriteYearToFile(doubleYear)
	}

	// Good years to test dragons: 2048 (for OUTWARD) and 2049 (for INWARD).

	// TODO check documentation. It's outdated in places
	// TODO English in the Czech text
	// TODO Add a --nechapu option which explains how our calendar works.
	// TODO Change from the document: The Pluto planet may repeat during Batman Days, because there could be up to 12 of the Days, not 9.
}
