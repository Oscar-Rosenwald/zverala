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
	doubleYear, kyear, targetDate, foundCached := command_line.RequestDate()
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

	if !foundCached {
		command_line.WriteYearToFile(doubleYear)
	}
}
