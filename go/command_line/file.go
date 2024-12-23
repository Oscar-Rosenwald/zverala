package command_line

import (
	"bufio"
	"os"
	"zverala/utils"
)

// WriteYearToFile stores year in its cache form in File.
func WriteYearToFile(year doubleYear) {
	if !SaveToFile {
		return
	}

	// Even if the file didn't exist when we called the app, cachedYear created
	// it, so it's safe now to open it without checking.
	utils.PrintDebug("Otevírám soubor %s abych zapsal dvojrok %s", File, year.ToString())

	f, err := os.OpenFile(File, os.O_APPEND|os.O_WRONLY, 0)
	utils.HandleError(err)
	defer f.Close()

	file := bufio.NewWriter(f)
	_, err = file.WriteString(year.ToCache())
	utils.HandleError(err)
	err = file.Flush()
	utils.HandleError(err)
	utils.PrintDebug("Dvojrok %s uložen", year.ToCache())
}
