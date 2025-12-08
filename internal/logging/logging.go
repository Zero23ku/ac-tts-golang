package logging

import (
	"os"
	"path/filepath"
	"time"
)

var logName = "error-"
var extension = ".txt"
var timeFormat = "02-01-2006-15:04:05"

func CreateLog(err error) {
	t := time.Now()
	tFormat := t.Format(timeFormat)
	fileName := logName + tFormat + extension

	path := filepath.Join(os.TempDir(), fileName)
	os.WriteFile(path, []byte(err.Error()), 0644)

}
