package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/Polshkrev/BollocksBot/models"
	"github.com/Polshkrev/BollocksBot/models/bot"
	"github.com/Polshkrev/BollocksBot/models/bot/command"
	"github.com/Polshkrev/BollocksBot/models/irc"
	"github.com/Polshkrev/BollocksBot/settings"
	"github.com/Polshkrev/BollocksBot/setup"
	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/fayl"
	"github.com/joho/godotenv"
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

// Write a variadic amount of messages through the given bot's IRC client.
func writeMessages(bot *bot.Bot, messages ...string) {
	var message string
	for _, message = range messages {
		bot.Write(message)
	}
}

// Check if the given key is stored in the system's enviornment.
// Returns true if the given key is set within the system's enviornment.
func check(key string) bool {
	var ok bool
	_, ok = os.LookupEnv(key)
	return ok
}

func main() {
	var settingsPath *fayl.Path = gopolutils.Must(setup.Configuration(settings.Path))

	var configuration settings.Settings = settings.Read(settingsPath)

	var token string

	if !check(configuration.Twitch.TokenKey) {
		var enviornmentFile *fayl.Path = gopolutils.Must(setup.Configuration(fayl.PathFrom(configuration.EnviornmentFilename)))
		var enviornment collections.Mapping[string, string] = loadEnviorment(enviornmentFile)
		token = *gopolutils.Must(enviornment.At(configuration.Twitch.TokenKey))
	} else {
		token = os.Getenv(configuration.Twitch.TokenKey)
	}

	var ircUrl url.URL = urlParse(fmt.Sprintf("%s://%s:%d", configuration.Twitch.Scheme, configuration.Twitch.BaseUrl, configuration.Twitch.Port))
	var authMessage string = fmt.Sprintf("%s %s", models.Authenticate, token)
	var nameMessage string = fmt.Sprintf("%s %s", models.Name, configuration.BotName)
	var joinMessage string = fmt.Sprintf("%s #%s", models.Join, configuration.Twitch.ChannelName)

	var irc *irc.IRC = irc.New(ircUrl.String())
	defer irc.Close()

	var bot *bot.Bot = bot.New(irc, configuration)
	writeMessages(bot, authMessage, nameMessage, joinMessage)

	command.Setup()

	var waitGroup *sync.WaitGroup = new(sync.WaitGroup)

	waitGroup.Go(func() {
		bot.Read()
	})

	var reader *bufio.Reader = bufio.NewReader(os.Stdin)
	waitGroup.Go(func() {
		bot.WaitForInput(reader)
	})

	waitGroup.Wait()
}
