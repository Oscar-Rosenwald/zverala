// -*- eval: (hs-minor-mode 1); -*-
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	ktime "zverala/klvanistic_time"
	"zverala/utils"
)

var saveToFile bool = true

// Handles command line argument
func parseArgs() {
	args := os.Args
	if len(args) == 1 {
		return // All default values are already defined
	}

	for i, arg := range args[1:] {
		switch arg {
		case "-h",
			"--help":
			printHelp()
			os.Exit(0)
		case "-n":
			saveToFile = false
		case "-f":
			if len(args) <= i+1 {
				fmt.Println("Po -f musí následovat jméno souboru")
				os.Exit(1)
			}
			file = args[i+1]
		case "--debug":
			utils.Debug = true
		case "--debug-debug":
			utils.DebugDebug = true
		default:
			if arg != file {
				fmt.Printf("Neznámý argument %s\n", arg)
				os.Exit(1)
			}
			// It is the filename from the -f option; do nothing
		}
	}
}

// requestYearInfo prompts for and reads from stdin information about the kyear
// in question.
func requestYearInfo() (doubleYear, kYear, bool) {
	stdinReader := bufio.NewReader(os.Stdin)

	readOption := func(prompt string) int {
		fmt.Print(prompt)
		ret, err := stdinReader.ReadString('\n')
		utils.HandleError(err)

		if strings.HasSuffix(ret, "\n") {
			ret = ret[:len(ret)-1]
		}

		retConv, err := strconv.Atoi(ret)
		utils.HandleError(err)
		return retConv
	}

	readSolstice := func(year int) time.Time {
		sol := readOption(fmt.Sprintf("Zimni slunovrat roku %d (den v prosinci): ", year))
		return time.Date(year, 12, sol, 11, 0, 0, 0, time.UTC)
	}

	utils.PrintInfo("Martinismus počítá čas ve dvojrocích. Dvojrok začíná o zimního slunovratu a nesedí tudíž přesně na standardní dataci.")
	utils.PrintInfo("")
	utils.PrintInfo("Tento program není schopen určit datum slunovratu. Připravte se prosím ho zadat.")
	utils.PrintInfo("")

	year := readOption("Rok prvního slunovratu: ")

	sol1, sol2, sol3, found := cachedYear(year)
	if found {
		var kyear kYear
		outKyear := ktime.ComputeKyear(sol1, sol2)
		inKyear := ktime.ComputeKyear(sol2, sol3)

		if sol2.Year() > year {
			kyear = outKyear
		} else {
			kyear = inKyear
		}

		return doubleYear{
			OutKyear: outKyear,
			InKyear:  inKyear,
			EndTime:  sol3,
			Length:   outKyear.Length + inKyear.Length,
		}, kyear, found
	}

	sol1 = readSolstice(year)
	sol2 = readSolstice(year + 1)
	kYear := ktime.ComputeKyear(sol1, sol2)

	switch kYear.Direction {
	case ktime.OUT:
		utils.PrintInfo("Krok v zadaném rozmezí je odstředný. Nyní potřebuji informace o následujícím kroku.")
		endSol := readSolstice(year + 2)
		inKyear := ktime.ComputeKyear(sol2, endSol)

		return doubleYear{
			OutKyear: kYear,
			InKyear:  inKyear,
			EndTime:  endSol,
			Length:   kYear.Length + inKyear.Length,
		}, kYear, found
	case ktime.IN:
		utils.PrintInfo("Krok v zadaném rozmezí je soustředný. Nyní potřebuji informace o předchozím kroku.")
		startSol := readSolstice(year - 1)
		outKyear := ktime.ComputeKyear(startSol, sol1)

		return doubleYear{
			OutKyear: outKyear,
			InKyear:  kYear,
			EndTime:  sol2,
			Length:   outKyear.Length + kYear.Length,
		}, kYear, found
	}

	utils.PrintDebug("Spracovávám krok %s", kYear.ToReadableString())
	utils.HandleError(fmt.Errorf("Cannot compute doubleyear with direction %c", kYear.Direction.ToChar()))
	return doubleYear{}, kYear, found
}

// Asks for current year and the winter solstices
func printHelp() {
	utils.PrintInfo("Vypočítej současný klvanistický rok. Všechny potřebné údaje budou automaticky vyžádány.")
	utils.PrintInfo("Program může zapsat do souboru a tisknout na konzoli, nebo jen tisknout.")
	utils.PrintInfo("Pokud rok už v souboru je, pouze tiskni. Defaultní chování je tisk a zápis do Zverala.txt.")
	utils.PrintInfo("-n            ... Pouze tisknout, nezapisovat")
	utils.PrintInfo("-f <soubor>   ... Hledat v / zapsat do souboru")
	utils.PrintInfo("-h --help     ... Zobrazit tento text")
	utils.PrintInfo("--debug       ... Tisknout extra informace")
	utils.PrintInfo("--debug-debug ... Tisknout extra EXTRA informace")
}

// printCreatures prints creatures with their day ranges in both normal and
// Klvanistic calendars.
func printCreatures(creatures []Creature, doubleYear doubleYear, kYear kYear, padToColumn, maxDaysLength int) {
	utils.PrintInfo("")
	var (
		currentDayRoller = kYear.NormalYearStart
		dateFormat       = "2.1. 2006"
	)

	for _, b := range creatures {
		lenB := utils.StrLen(string(b.name))
		dotPadB := utils.Pad(padToColumn, lenB, ".")
		days := b.days
		lenDays := utils.DaysDigits(days)
		spacePad := utils.Pad(maxDaysLength, lenDays, " ")

		if days == 0 {
			fmt.Printf("%s%s%d%sdnů\n",
				b.name,
				dotPadB,
				days,
				spacePad)
			continue
		}

		first := currentDayRoller
		// Add days-1, because the last day belongs to the next creature, unless
		// this is the last creature of the kyear, which is not followed by
		// anything.
		currentDayRoller = currentDayRoller.AddDate(0, 0, days-1)
		last := currentDayRoller

		// Now add the remaining 1 so the next creature starts on their day.
		currentDayRoller = currentDayRoller.AddDate(0, 0, 1)

		kStart := ktime.TimeToKlvanisticDate(first, doubleYear)
		kEnd := ktime.TimeToKlvanisticDate(last, doubleYear)

		fmt.Printf("%s%s%d%s%s: % 11s - %-11s | %s - %-11s\n",
			b.name,
			dotPadB,
			days,
			spacePad,
			utils.DayString(days),
			first.Format(dateFormat),
			last.Format(dateFormat),
			kStart.ToString(),
			kEnd.ToString(),
		)
	}

}
