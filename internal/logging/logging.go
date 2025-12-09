package logging

import (
	"os"
	"path/filepath"
	"time"
)

var logName = "error-"
var extension = ".txt"
var timeFormat = "02-01-2006-15:04:05"
var newLine = "\n"

func CreateLog(place string, err error) {
	t := time.Now()
	tFormat := t.Format(timeFormat)
	fileName := logName + tFormat + extension

	path := filepath.Join(os.TempDir(), fileName)
	file, openErr := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if openErr != nil {
		//Todo se derrumbó, dentro de mi, dentro de mi
		panic(openErr)
	}
	defer file.Close()

	log := tFormat + "-" + place + newLine + err.Error() + newLine
	if _, writeErr := file.WriteString(log); writeErr != nil {
		//WHAT IS HAPPENING AAAAA
		panic(writeErr)
	}

}
