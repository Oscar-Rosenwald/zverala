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

var Debug = false

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
		default:
			if arg != file {
				fmt.Printf("Neznámý argument %s\n", arg)
				os.Exit(1)
			}
			// It is the filename from the -f option; do nothing
		}
	}
}

// Asks for current year and the winter solstices
func requestYearInfo() (normalYear int, lastSolstice, nextSolstice time.Time) {
	// TODO Remove this block - I​ got tired of inputting the same year over and over again.
	if false {
		fmt.Printf("RUNNING IN TEST MODE - NO YEAR NEEDS TO BE GIVEN\n\n") // ___
		return 2021, time.Date(2020, time.December, 21, 10, 10, 10, 0, time.UTC), time.Date(2021, time.December, 21, 10, 10, 10, 0, time.UTC)
	}

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

	year := readOption("Tento rok: ")
	solstice1 := readOption(fmt.Sprintf("Zimni slunovrat roku %d (den v prosinci): ", year-1))
	solstice2 := readOption(fmt.Sprintf("Zimni slunovrat roku %d (den v prosinci): ", year))

	sol1 := time.Date(year-1, 12, solstice1, 12, 0, 0, 0, time.Local)
	sol2 := time.Date(year, 12, solstice2, 12, 0, 0, 0, time.Local)

	fmt.Printf("Propočítávám rok od %s do %s\n", sol1, sol2)

	return year, sol1, sol2
}

func printHelp() {
	fmt.Println("Vypočítej současný klvanistický rok. Všechny potřebné údaje budou automaticky vyžádány.")
	fmt.Println("Program může zapsat do souboru a tisknout na konzoli, nebo jen tisknout.")
	fmt.Println("Pokud rok už v souboru je, pouze tiskni. Defaultní chování je tisk a zápis do Zverala.txt.")
	fmt.Println("-n          ... Pouze tisknout, nezapisovat")
	fmt.Println("-f <soubor> ... Hledat v / zapsat do souboru")
	fmt.Println("-h --help   ... Zobrazit tento text")
}

func readYearFromFile(doubleYear kYear) (yearDetail string, success bool) {
	printDebug("Otevírám soubor %s", file)
	f, err := os.Open(file)
	handleError(err)
	defer f.Close()
	file := bufio.NewScanner(f)

	var noEnd = true
	for noEnd {
		noEnd = file.Scan()
		handleError(file.Err())

		line := file.Text()
		fmt.Printf("Read line '%s'\n", line) // ___

		printDebug("Porovnávám s uloženým krokem %s", doubleYear.toReadableString())
		if line == doubleYear.toReadableString() {
			printDebug("Našel jsem shodu v krocích")
			for noEnd {
				yearDetail += fmt.Sprintf("\n%s", line)
				noEnd = file.Scan()
				handleError(file.Err())
				line = file.Text()
				if line == "" || !noEnd {
					return yearDetail, true
				}
			}
			fmt.Printf("No more lines to read") // ___
			return "", false
		}
	}

	fmt.Printf("End of file\n") // ___
	printDebug("Krok není uložen v souboru")
	return "", false
}

func printCreatures(creatures []Creature, doubleYear kYear, padToColumn, maxDaysLength int) {
	var (
		currentDayRoller = doubleYear.normalYearStart
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
		// Add days-1, because the last day belongs to the next creature.
		currentDayRoller = currentDayRoller.AddDate(0, 0, days-1)
		last := currentDayRoller
		// Now add the remaining 1 so the next creature starts on their day.
		currentDayRoller = currentDayRoller.AddDate(0, 0, 1)

		fmt.Printf("%s%s%d%s%s: % 11s - % 11s\n",
			b.name,
			dotPadB,
			days,
			spacePad,
			dayString(days),
			first.Format(dateFormat),
			last.Format(dateFormat))
	}

}
