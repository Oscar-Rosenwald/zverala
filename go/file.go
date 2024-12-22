package main

import (
	"bufio"
	"os"
	"zverala/command_line"
	"zverala/utils"
)

func writeYearToFile(year doubleYear) {
	if !command_line.SaveToFile {
		return
	}

	// Even if the file didn't exist when we called the app, cachedYear created
	// it, so it's safe now to open it without checking.
	utils.PrintDebug("Otevírám soubor %s", command_line.File)

	f, err := os.OpenFile(command_line.File, os.O_APPEND|os.O_WRONLY, 0)
	utils.HandleError(err)
	defer f.Close()

	file := bufio.NewWriter(f)
	_, err = file.WriteString(year.ToCache())
	utils.HandleError(err)
	err = file.Flush()
	utils.HandleError(err)
	utils.PrintDebug("Dvojrok %s uložen", year.ToCache())
}
