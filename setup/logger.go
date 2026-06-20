package setup

import (
	"fmt"
	"time"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/fayl"
)

const (
	suffix fayl.Suffix = fayl.Log
)

// Obtain the current formatted time.
// Returns a [fayl.Path] formatted with the current time.
func getLogFile(format string) *fayl.Path {
	var logTime string = time.Now().Format(format)
	var logSuffix string = gopolutils.Must(fayl.StringFromSuffix(suffix))
	var logFile string = fmt.Sprintf("%s.%s", logTime, logSuffix)
	return fayl.PathFrom(logFile)
}

// Setup the log file based on the given folder and time format.
// Returns the log file path.
// If the given log file path can not be created, the associated error is returned.
func Logging(logFolder, format string) (*fayl.Path, *gopolutils.Exception) {
	var configurationPath *fayl.Path = makeConfiguration()
	var logPath *fayl.Path = configurationPath.JoinAs(logFolder)
	var logEntry *fayl.Entry = fayl.NewEntry(logPath)
	logEntry.SetType(fayl.DirectoryType)
	var except *gopolutils.Exception = makeEntry(logEntry)
	if except != nil {
		return nil, except
	}
	var logFile *fayl.Path = getLogFile(format)
	var finalPath *fayl.Path = logPath.Join(*logFile)
	var finalEntry *fayl.Entry = fayl.NewEntry(finalPath)
	if finalPath.Exists() {
		return finalPath, nil
	}
	except = makeEntry(finalEntry)
	if except != nil {
		return nil, except
	}
	return finalPath, nil
}
