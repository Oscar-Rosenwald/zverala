package main

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO Change for the appropriate file. Make it a relative path, possibly configurable by the user.
var file = "/Users/cyrilsaroch/Documents/Martinismus/Zverala/Zverala2.txt"

func writeYearToFile(year doubleYear) {
	if !saveToFile {
		return
	}

	// Even if the file didn't exist when we called the app, cachedYear created
	// it, so it's safe now to open it without checking.
	printDebug("Otevírám soubor %s", file)

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0)
	handleError(err)
	defer f.Close()

	file := bufio.NewWriter(f)
	_, err = file.WriteString(year.toCache())
	handleError(err)
	err = file.Flush()
	handleError(err)
	printDebug("Dvojrok %s uložen", year.toCache())
}

func cachedYear(year int) (sol1, sol2, sol3 time.Time, found bool) {
	if !saveToFile {
		return time.Time{}, time.Time{}, time.Time{}, false
	}

	var f *os.File

	_, err := os.Stat(file)
	if err == nil {
		printDebug("Otevírám soubor %s", file)
		f, err = os.Open(file)
	} else if errors.Is(err, fs.ErrNotExist) {
		printDebug("Vytvářím soubor %s", file)
		f, err = os.Create(file)
	}

	handleError(err)
	defer f.Close()

	file := bufio.NewScanner(f)
	for {
		noEnd := file.Scan()
		handleError(file.Err())
		if !noEnd {
			break
		}
		line := file.Text()
		if line == "" {
			break
		}

		printDEBUG("Porovnávám dvojrok %d s řádkem %s", year, line)
		parts := strings.Split(line, ":")
		if len(parts) != 5 {
			printInfo("EROR: SOUBOR %d NEMÁ SPRÁVNÝ FORMÁT", f.Name())
			break
		}

		firstYear, err := strconv.Atoi(parts[1])
		handleError(err)
		midDay, err := strconv.Atoi(parts[3])
		handleError(err)
		midSol := time.Date(firstYear+1, 12, midDay, 0, 0, 0, 0, time.Local)

		printDEBUG("První rok: %d, prostřední slunovrat: %s", firstYear, midSol.String())

		if year == firstYear || year == firstYear+1 {
			firstDay, err := strconv.Atoi(parts[2])
			handleError(err)
			solStart := time.Date(firstYear, 12, firstDay, 0, 0, 0, 0, time.Local)

			lastDay, err := strconv.Atoi(parts[4])
			handleError(err)
			solEnd := time.Date(firstYear+2, 12, lastDay, 0, 0, 0, 0, time.Local)
			return solStart, midSol, solEnd, true
		}
	}

	return time.Time{}, time.Time{}, time.Time{}, false
}
