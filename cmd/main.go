package main

import (
	"fmt"
	"net/url"

	"github.com/Polshkrev/BollocksBot/models"
	"github.com/Polshkrev/BollocksBot/models/irc"
	"github.com/Polshkrev/BollocksBot/settings"
	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/fayl"
	"github.com/joho/godotenv"
)

var (
	config      settings.Settings                   = settings.Read(settings.Path)                                                                         // Global configuration of the application.
	enviornment collections.Mapping[string, string] = loadEnviorment(fayl.PathFrom(config.EnviornmentFilename))                                            // Global enviornment.
	token       string                              = *gopolutils.Must(enviornment.At(config.Twitch.TokenKey))                                             // Twitch token.
	ircUrl      url.URL                             = urlParse(fmt.Sprintf("%s://%s:%d", config.Twitch.Scheme, config.Twitch.BaseUrl, config.Twitch.Port)) // Twitch IRC url.
	nameMessage string                              = fmt.Sprintf("%s %s", models.Name, config.BotName)                                                    // Login message.
	authMessage string                              = fmt.Sprintf("%s %s", models.Authenticate, token)                                                     // Token message.
	joinMessage string                              = fmt.Sprintf("%s #%s", models.Join, config.Twitch.ChannelName)                                        // Join message.
)

// Load a given .env file as a [collections.Mapping].
// Returns a [collections.Mapping] of a given .env file.
// If the given file does not exist, this function `panics` with a [gopolutils.FileNotFoundError].
// If the .env file can not be read, this function `panics` with an [gopolutils.OSError].
// If the key is already in the mapping, instead of just quietly not inserting into the mapping, this function `panics` with a [gopolutils.KeyError].
func loadEnviorment(file *fayl.Path) collections.Mapping[string, string] {
	if !file.Exists() {
		panic(gopolutils.NewNamedException(gopolutils.FileNotFoundError, "'%s' does not exist.", file))
	}
	var result collections.Mapping[string, string] = collections.NewMap[string, string]()
	var raw map[string]string
	var readError error
	raw, readError = godotenv.Read(file.String())
	if readError != nil {
		panic(gopolutils.NewNamedException(gopolutils.OSError, "%s\n", readError.Error()))
	}
	var key, value string
	for key, value = range raw {
		var insertExcept *gopolutils.Exception = result.Insert(key, value)
		if insertExcept != nil {
			panic(insertExcept)
		}
	}
	return result
}

// Parse a given raw URL.
// Returns a [url.URL] representation of the given raw URL.
// If the raw URL can not be parsed, this function `panics` with a [gopolutils.RuntimeError].
func urlParse(raw string) url.URL {
	var result *url.URL
	var parseError error
	result, parseError = url.Parse(raw)
	if parseError != nil {
		panic(gopolutils.NewNamedException(gopolutils.RuntimeError, "%s\n", parseError.Error()))
	}
	return *result
}

func main() {
	var irc *irc.IRC = irc.New(ircUrl.String())
	defer irc.Close()
	var except *gopolutils.Exception = irc.Write(authMessage)
	if except != nil {
		panic(except)
	}
	except = irc.Write(nameMessage)
	if except != nil {
		panic(except)
	}
	except = irc.Write(joinMessage)
	if except != nil {
		panic(except)
	}
	for {
		fmt.Print(gopolutils.Must(irc.Read()))
	}
}
