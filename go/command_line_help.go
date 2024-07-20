// -*- eval: (hs-minor-mode 1); -*-
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Print helpful debug logs.
var Debug = false

// Print unhelpful debug logs.
var DebugDebug = false

func printDebug(msg string, args ...interface{}) {
	if !Debug {
		return
	}

	if len(args) > 0 {
		fmt.Printf(msg+"\n", args...)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

func printDEBUG(msg string, args ...interface{}) {
	if !DebugDebug {
		return
	}

	if len(args) > 0 {
		fmt.Printf(msg+"\n", args...)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

func printInfo(msg string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf(msg+"\n", args...)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

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
			Debug = true
		case "--debug-debug":
			DebugDebug = true
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
func requestYearInfo() (doubleYear, kYear) {
	stdinReader := bufio.NewReader(os.Stdin)

	readOption := func(prompt string) int {
		fmt.Print(prompt)
		ret, err := stdinReader.ReadString('\n')
		handleError(err)

		if strings.HasSuffix(ret, "\n") {
			ret = ret[:len(ret)-1]
		}

		retConv, err := strconv.Atoi(ret)
		handleError(err)
		return retConv
	}

	readSolstice := func(year int) time.Time {
		sol := readOption(fmt.Sprintf("Zimni slunovrat roku %d (den v prosinci): ", year))
		return time.Date(year, 12, sol, 11, 0, 0, 0, time.UTC)
	}

	printInfo("Martinismus počítá čas ve dvojrocích. Dvojrok začíná o zimního slunovratu a nesedí tudíž přesně na standardní dataci.")
	printInfo("")
	printInfo("Tento program není schopen určit datum slunovratu. Připravte se prosím ho zadat.")
	printInfo("")

	year := readOption("Začátek dvojroku (rok, ve kterém je první zimní slunovrat kroku): ")

	sol1, sol2, sol3, found := cachedYear(year)
	if found {
		var kyear kYear
		outKyear := computeKyear(sol1, sol2)
		inKyear := computeKyear(sol2, sol3)

		if sol2.Year() > year {
			kyear = outKyear
		} else {
			kyear = inKyear
		}

		return doubleYear{
			outKyear: outKyear,
			inKyear:  inKyear,
			endTime:  sol3,
			length:   outKyear.length + inKyear.length,
		}, kyear
	}

	sol1 = readSolstice(year)
	sol2 = readSolstice(year + 1)
	kYear := computeKyear(sol1, sol2)

	switch kYear.direction {
	case OUT:
		printInfo("Krok v zadaném rozmezí je odstředný. Nyní potřebuji informace o následujícím kroku.")
		endSol := readSolstice(year + 2)
		inKyear := computeKyear(sol2, endSol)

		return doubleYear{
			outKyear: kYear,
			inKyear:  inKyear,
			endTime:  endSol,
			length:   kYear.length + inKyear.length,
		}, kYear
	case IN:
		printInfo("Krok v zadaném rozmezí je soustředný. Nyní potřebuji informace o předchozím kroku.")
		startSol := readSolstice(year - 1)
		outKyear := computeKyear(startSol, sol1)

		return doubleYear{
			outKyear: outKyear,
			inKyear:  kYear,
			endTime:  sol2,
			length:   outKyear.length + kYear.length,
		}, kYear
	}

	printDebug("Spracovávám krok %s", kYear.toReadableString())
	handleError(fmt.Errorf("Cannot compute doubleyear with direction %c", kYear.direction.toChar()))
	return doubleYear{}, kYear
}

// Asks for current year and the winter solstices
func printHelp() {
	printInfo("Vypočítej současný klvanistický rok. Všechny potřebné údaje budou automaticky vyžádány.")
	printInfo("Program může zapsat do souboru a tisknout na konzoli, nebo jen tisknout.")
	printInfo("Pokud rok už v souboru je, pouze tiskni. Defaultní chování je tisk a zápis do Zverala.txt.")
	printInfo("-n            ... Pouze tisknout, nezapisovat")
	printInfo("-f <soubor>   ... Hledat v / zapsat do souboru")
	printInfo("-h --help     ... Zobrazit tento text")
	printInfo("--debug       ... Tisknout extra informace")
	printInfo("--debug-debug ... Tisknout extra EXTRA informace")
}

// printCreatures prints creatures with their day ranges in both normal and
// Klvanistic calendars.
func printCreatures(creatures []Creature, doubleYear doubleYear, kYear kYear, padToColumn, maxDaysLength int) {
	printInfo("")
	var (
		currentDayRoller = kYear.normalYearStart
		dateFormat       = "2.1. 2006"
	)

	for _, b := range creatures {
		lenB := strLen(string(b.name))
		dotPadB := pad(padToColumn, lenB, ".")
		days := b.days
		lenDays := daysDigits(days)
		spacePad := pad(maxDaysLength, lenDays, " ")

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

		kStart := timeToKlvanisticDate(first, doubleYear)
		kEnd := timeToKlvanisticDate(last, doubleYear)

		fmt.Printf("%s%s%d%s%s: % 11s - %-11s | %s - %-11s\n",
			b.name,
			dotPadB,
			days,
			spacePad,
			dayString(days),
			first.Format(dateFormat),
			last.Format(dateFormat),
			kStart.toString(),
			kEnd.toString(),
		)
	}

}
