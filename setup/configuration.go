package setup

import (
	"os"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/fayl"
	"github.com/Polshkrev/goserialize"
)

var (
	folder          string             = "BollocksBot"       // Configuration folder name. This should be set to the name of the bot.
	defaultSettings goserialize.Object = goserialize.Object{ // Default settings to copy to the configuration file.
		"botName":             "BollocksBot",
		"enviornmentFilename": ".env",
		"twitch": goserialize.Object{
			"tokenKey":    "TWITCH_TOKEN",
			"scheme":      "wss",
			"baseUrl":     "irc-ws.chat.twitch.tv",
			"port":        443,
			"channelName": "polshkrev",
			"today": goserialize.Object{
				"topic": "Making the chat bot.",
			},
		},
		"logging": goserialize.Object{
			"folder": "logs",
			"format": "2006-01",
		},
	}
)

// Obtain the configuration folder name.
// Returns the configuration path.
// If the configuration path can not be obtained, a [gopolutils.FileNotFound] error is returned.
func GetConfigurationPath() (*fayl.Path, *gopolutils.Exception) {
	var result string
	var configurationError error
	result, configurationError = os.UserConfigDir()
	if configurationError != nil {
		return nil, gopolutils.NewNamedException(gopolutils.FileNotFoundError, "%s", configurationError.Error())
	}
	return fayl.PathFrom(result), nil
}

// Create the given entry path.
// If the entry path can not be created, its assocciated error is returned.
func makeEntry(entry *fayl.Entry) *gopolutils.Exception {
	if entry.Path().Exists() {
		return nil
	}
	var parent *fayl.Path = gopolutils.Must(entry.Path().Parent())
	var parentEntry *fayl.Entry = fayl.NewEntry(parent)
	parentEntry.SetType(fayl.DirectoryType)
	if !parent.Exists() {
		return parentEntry.MakeDirectory()
	}
	return entry.Create()
}

// Setup the configuration based on the given path.
// Returns the configuration path.
// If the given configuration path can not be created, the associated error is returned.
func Configuration(path *fayl.Path) (*fayl.Path, *gopolutils.Exception) {
	var configurationPath *fayl.Path
	var except *gopolutils.Exception
	configurationPath, except = GetConfigurationPath()
	if except != nil {
		var configurationEntry *fayl.Entry = fayl.NewEntry(configurationPath)
		configurationEntry.SetType(fayl.DirectoryType)
		return nil, makeEntry(configurationEntry)
	}
	var botPath *fayl.Path = configurationPath.Join(*fayl.PathFrom(folder))
	var botEntry *fayl.Entry = fayl.NewEntry(botPath)
	botEntry.SetType(fayl.DirectoryType)
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
