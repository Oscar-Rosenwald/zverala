package main

import (
	"fmt"
	"zverala/command_line"
	ktime "zverala/klvanistic_time"
	"zverala/spiral"
	"zverala/utils"
)

func main() {
	command_line.ParseArgs()
	doubleYear, kyear, targetDate := command_line.RequestDate()
	utils.PrintDebug("Zpracovávám datum %s v kroce %s dvojroku %s", targetDate.String(), kyear.ToReadableString(), doubleYear.ToString())

	// In this case it's okay to ignore the fact that not all days are 24 hours,
	// because the difference will be in the order of hours, and we don't care
	// about that.
	daysFromKyearStart := int(targetDate.Sub(kyear.NormalYearStart).Hours()/24) + 1
	targetKDate := ktime.TimeToKlvanisticDate(targetDate, doubleYear)
	creaturesInOrder, _ := spiral.ComputeZverala(&kyear, &doubleYear)
	utils.PrintDebug("Datum %s je %d dní od začátku kroku %s", targetDate.String(), daysFromKyearStart, kyear.ToString())

	var creatureAtDate utils.Creature
	dayRoller := 0
	for _, creature := range creaturesInOrder {
		utils.PrintDebug("Patří datum bytosti %s (která působí po %d dní)", creature.Name, creature.Days)
		dayRoller += creature.Days
		if dayRoller >= int(daysFromKyearStart) {
			creatureAtDate = creature
			break
		}
	}

	if creatureAtDate.Name == "" {
		utils.HandleError(fmt.Errorf("Nemohu nalézt žádnou bytost %d dní po začátku kroku %s", daysFromKyearStart, kyear.ToReadableString()))
	}

	utils.PrintInfo("%s: %s", targetKDate.ToString(), creatureAtDate.Name)

	// TODO If the date is in the third year of a doubleyear (as in January to
	// December of the last year), then CachedYear doesn't consider it cached
	// even though it is. This makes sense, because CachedYear only looks at the
	// first two years, as the third solstice belongs into the next doubleyear,
	// but in this case where we are supplying a precise date we do actually
	// need it.

	// TODO After CachedYear can differentiate between 1st January 2049 and 31st
	// December 2049 as two different double years, we will want to cache any
	// new years. As of writing this we cannot do that because 2049 is a year
	// that according to CachedYear doesn't belong to the doubleyear starting
	// with 2047, but it DOES belong there until December.
}
