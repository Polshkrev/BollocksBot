package setup

import (
	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/fayl"
)

// Concurrently create a given entry.
func makeConcurrentEntry(entry *fayl.Entry, errorChannel chan<- *gopolutils.Exception) {
	defer close(errorChannel)
	if entry.Path().Exists() {
		errorChannel <- nil
		return
	}
	var parent *fayl.Path = gopolutils.Must(entry.Path().Parent())
	var parentEntry *fayl.Entry = fayl.NewEntry(parent)
	parentEntry.SetType(fayl.DirectoryType)
	if !parent.Exists() {
		errorChannel <- parentEntry.MakeDirectory()
		return
	}
	errorChannel <- entry.Create()
}

// Create the given entry path.
// If the entry path can not be created, its assocciated error is returned.
func makeEntry(entry *fayl.Entry) *gopolutils.Exception {
	var errorChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception, 1)
	go makeConcurrentEntry(entry, errorChannel)
	var except *gopolutils.Exception = <-errorChannel
	return except
}

// Create the configuration folder based on a given path.
func makeConcurrentConfiguration(resultChannel chan<- *fayl.Path) {
	defer close(resultChannel)
	var configurationPath *fayl.Path = fayl.Configuration()
	resultChannel <- configurationPath
}

// Create the configuration folder if it does not exist.
// Retruns the [fayl.Path] of the configuration folder.
func makeConfiguration() *fayl.Path {
	var resultChannel chan *fayl.Path = make(chan *fayl.Path, 1)
	go makeConcurrentConfiguration(resultChannel)
	var result *fayl.Path = <-resultChannel
	return result
}

// Setup the configuration based on the given path.
// Returns the configuration path.
// If the given configuration path can not be created, the associated error is returned.
func Configuration(path *fayl.Path) (*fayl.Path, *gopolutils.Exception) {
	var configurationPath *fayl.Path = makeConfiguration()

	var userFolderPath *fayl.Path = configurationPath.JoinAs(userFolder)
	var userFolderEntry *fayl.Entry = fayl.NewEntry(userFolderPath)

	var except *gopolutils.Exception = makeEntry(userFolderEntry)
	if except != nil {
		return nil, except
	}

	var botPath *fayl.Path = userFolderPath.Join(*fayl.PathFrom(folder))
	var botEntry *fayl.Entry = fayl.NewEntry(botPath)

	except = makeEntry(botEntry)
	if except != nil {
		return nil, except
	}

	var settingsPath *fayl.Path = botPath.Join(*path)
	var settingsEntry *fayl.Entry = fayl.NewEntry(settingsPath)
	if settingsPath.Exists() {
		return settingsPath, nil
	}

	except = makeEntry(settingsEntry)
	if except != nil {
		return nil, except
	}

	except = fayl.WriteObject(settingsPath, &defaultSettings)
	if except != nil {
		return nil, except
	}

	return settingsPath, nil
}
