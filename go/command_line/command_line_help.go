// -*- eval: (hs-minor-mode 1); -*-
package command_line

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
	ktime "zverala/klvanistic_time"
	"zverala/utils"
)

type doubleYear = ktime.DoubleYear
type kYear = ktime.KYear

var (
	SaveToFile bool = true
	File            = "../Zverala-cache.txt"
)

// Handles command line argument
func ParseArgs() {
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
			SaveToFile = false
		case "-f":
			if len(args) <= i+1 {
				fmt.Println("Po -f musí následovat jméno souboru")
				os.Exit(1)
			}
			File = args[i+1]
		case "--debug":
			utils.Debug = true
		case "--debug-debug":
			utils.DebugDebug = true
		default:
			if arg != File {
				fmt.Printf("Neznámý argument %s\n", arg)
				os.Exit(1)
			}
			// It is the filename from the -f option; do nothing
		}
	}
}

// requestYearInfo prompts for and reads from stdin information about the kyear
// in question. It is to be used to print the whole kyear.
func RequestYearInfo() (doubleYear, kYear, bool) {
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

	sol1, sol2, sol3, found := CachedYear(year, false, 0)
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
	utils.HandleError(fmt.Errorf("Nemohu zpracovat dvojrok ve směru %c", kYear.Direction.ToChar()))
	return doubleYear{}, kYear, found
}

// TODO tests
//
// RequestDate prompts the user for info about a date in question, computing the
// doublyear surrounding it. The date inputed is returned in targetDate.
func RequestDate() (doubleyear doubleYear, kyear kYear, targetDate time.Time, foundCached bool) {
	stdinReader := bufio.NewReader(os.Stdin)

	readOption := func(prompt string) int {
		fmt.Print(prompt + " ")
		ret, err := stdinReader.ReadString('\n')
		utils.HandleError(err)

		if strings.HasSuffix(ret, "\n") {
			ret = ret[:len(ret)-1]
		}

		retConv, err := strconv.Atoi(ret)
		utils.HandleError(err)
		return retConv
	}

	year := readOption("Který rok vás zajímá?")
	month := readOption("Který měsíc vás zajíma?")
	day := readOption("Který den vás zajíma?")

	// Set the hour of the target date to be 12 so it's always guaranteed to be
	// after the cached solstice, which is always at midnight, and after the
	// reference date's timestamp, which is 11:00 UTC.
	targetDate = time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
	startIndex := 0

	for {
		sol1, sol2, sol3, found := CachedYear(year, true, startIndex)
		if !found {
			break
		}

		if sol1.Before(targetDate) && sol3.After(targetDate) {
			outKyear := ktime.ComputeKyear(sol1, sol2)
			inKyear := ktime.ComputeKyear(sol2, sol3)

			currentKyear := outKyear
			if targetDate.After(sol2) {
				currentKyear = inKyear
			}

			return doubleYear{
				OutKyear: outKyear,
				InKyear:  inKyear,
				EndTime:  sol3,
				Length:   outKyear.Length + inKyear.Length,
			}, currentKyear, targetDate, true
		}

		startIndex++
	}

	// We need to construct the doubleyear. First we need to find out if we're
	// before or after the solstice. If we're after it, we can ask for the
	// solstice of the request year. If we're before it, we need to ask for the
	// solistice of the previous year.

	readSolstice := func(year int) time.Time {
		sol := readOption(fmt.Sprintf("Zimni slunovrat roku %d (den v prosinci):", year))
		return time.Date(year, 12, sol, 11, 0, 0, 0, time.UTC)
	}

	solSameYear := readSolstice(year)
	var currnetKyear kYear
	var solNextYear, solPreviousYear time.Time

	if targetDate.After(solSameYear) {
		solNextYear = readSolstice(year + 1)
		currnetKyear = ktime.ComputeKyear(solSameYear, solNextYear)
	} else {
		solPreviousYear = readSolstice(year - 1)
		currnetKyear = ktime.ComputeKyear(solPreviousYear, solSameYear)
	}

	switch currnetKyear.Direction {
	case ktime.OUT:
		utils.PrintInfo("Krok v zadaném rozmezí je odstředný. Nyní potřebuji informace o následujícím kroku.")
		endSol := readSolstice(year + 1)
		inKyear := ktime.ComputeKyear(solNextYear, endSol)

		return doubleYear{
			OutKyear: currnetKyear,
			InKyear:  inKyear,
			EndTime:  endSol,
			Length:   currnetKyear.Length + inKyear.Length,
		}, currnetKyear, targetDate, false
	case ktime.IN:
		utils.PrintInfo("Krok v zadaném rozmezí je soustředný. Nyní potřebuji informace o předchozím kroku.")
		startSol := readSolstice(year - 2)
		outKyear := ktime.ComputeKyear(startSol, solPreviousYear)

		return doubleYear{
			OutKyear: outKyear,
			InKyear:  currnetKyear,
			EndTime:  solSameYear,
			Length:   outKyear.Length + currnetKyear.Length,
		}, currnetKyear, targetDate, false
	}

	utils.PrintDebug("Spracovávám krok %s", currnetKyear.ToReadableString())
	utils.HandleError(fmt.Errorf("Nemohu zpracovat dvojrok ve směru %c", currnetKyear.Direction.ToChar()))
	return doubleYear{}, kYear{}, time.Time{}, false
}

func PrintCreatures(creatures []utils.Creature, doubleYear doubleYear, kYear kYear, padToColumn, maxDaysLength int) {
	utils.PrintInfo("")
	var (
		currentDayRoller = kYear.NormalYearStart
		dateFormat       = "2.1. 2006"
	)

	for _, b := range creatures {
		lenB := utils.StrLen(string(b.Name))
		dotPadB := utils.Pad(padToColumn, lenB, ".")
		days := b.Days
		lenDays := utils.DaysDigits(days)
		spacePad := utils.Pad(maxDaysLength, lenDays, " ")

		if days == 0 {
			fmt.Printf("%s%s%d%sdnů\n",
				b.Name,
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
			b.Name,
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

// CachedYear finds year in File, taking into account that the file caches whole
// doubleyears. It returns the three solstices of the doubleyear (start, middle,
// and end) along with a flag that says whether a match was found.
//
// CachedYear is targetted at the zverala app, which only cares about whole
// years. If you want to use it in the calendar app, a closer examination of the
// dates must be undergone. To enable this, set checkSurroundingYears to true.
//
// The first startIndex matched lines are ignored. This way you may loop over
// this function multiple times, incrementing startIndex, and get different
// results every time.
func CachedYear(year int, checkSurroundingYears bool, startIndex int) (sol1, sol2, sol3 time.Time, found bool) {
	if !SaveToFile {
		return time.Time{}, time.Time{}, time.Time{}, false
	}

	var f *os.File

	_, err := os.Stat(File)
	if err == nil {
		utils.PrintDebug("Otevírám soubor %s", File)
		f, err = os.Open(File)
	} else if errors.Is(err, fs.ErrNotExist) {
		utils.PrintDebug("Vytvářím soubor %s", File)
		f, err = os.Create(File)
	}

	utils.HandleError(err)
	defer f.Close()

	file := bufio.NewScanner(f)
	for {
		noEnd := file.Scan()
		utils.HandleError(file.Err())
		if !noEnd {
			break
		}
		line := file.Text()
		if line == "" {
			break
		}

		if startIndex > 0 {
			startIndex--
			continue
		}

		utils.PrintDEBUG("Porovnávám dvojrok %d s řádkem %s", year, line)
		parts := strings.Split(line, ":")
		if len(parts) != 5 {
			utils.PrintInfo("EROR: SOUBOR %d NEMÁ SPRÁVNÝ FORMÁT", f.Name())
			break
		}

		firstYear, err := strconv.Atoi(parts[1])
		utils.HandleError(err)
		midDay, err := strconv.Atoi(parts[3])
		utils.HandleError(err)
		midSol := time.Date(firstYear+1, 12, midDay, 0, 0, 0, 0, time.Local)

		utils.PrintDEBUG("První rok: %d, prostřední slunovrat: %s", firstYear, midSol.String())

		if year == firstYear || year == firstYear+1 || (checkSurroundingYears && year == firstYear+2) {
			firstDay, err := strconv.Atoi(parts[2])
			utils.HandleError(err)
			solStart := time.Date(firstYear, 12, firstDay, 0, 0, 0, 0, time.Local)

			lastDay, err := strconv.Atoi(parts[4])
			utils.HandleError(err)
			solEnd := time.Date(firstYear+2, 12, lastDay, 0, 0, 0, 0, time.Local)
			return solStart, midSol, solEnd, true
		}
	}

	return time.Time{}, time.Time{}, time.Time{}, false
}
