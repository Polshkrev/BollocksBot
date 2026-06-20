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

// Setup the configuration based on the given path.
// Returns the configuration path.
// If the given configuration path can not be created, the associated error is returned.
func Logging(logFolder, format string) (*fayl.Path, *gopolutils.Exception) {
	var configurationPath *fayl.Path = makeConfiguration()

	var botPath *fayl.Path = configurationPath.Join(*fayl.PathFrom(folder))
	var botEntry *fayl.Entry = fayl.NewEntry(botPath)
	botEntry.SetType(fayl.DirectoryType)

	var except *gopolutils.Exception = makeEntry(botEntry)
	if except != nil {
		return nil, except
	}
	var logFolderPath *fayl.Path = botPath.JoinAs(logFolder)
	var logPath *fayl.Path = logFolderPath.Join(*getLogFile(format))
	if logFolderPath.Exists() {
		return logPath, nil
	}

	var logEntry *fayl.Entry = fayl.NewEntry(logPath)
	except = makeEntry(logEntry)
	if except != nil {
		return nil, except
	}

	return logPath, nil
}
