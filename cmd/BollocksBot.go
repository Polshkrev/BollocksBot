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
	"github.com/Polshkrev/BollocksBot/models/enviornment"
	"github.com/Polshkrev/BollocksBot/models/irc"
	"github.com/Polshkrev/BollocksBot/settings"
	"github.com/Polshkrev/BollocksBot/setup"
	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/fayl"
	"github.com/Polshkrev/goserialize"
)

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
	var i int
	for i = range messages {
		bot.Write(messages[i])
	}
}

// Set the today command's enviornment variable.
func setToday(value string) {
	var except *gopolutils.Exception = enviornment.Set(command.Today.String(), value)
	if except != nil {
		panic(except)
	}
}

// Set the links within the commands's enviornment.
func setLinks(links goserialize.Object) {
	var key string
	var value any
	for key, value = range links {
		var except *gopolutils.Exception = enviornment.Set(key, value.(string))
		if except != nil {
			panic(except)
		}
	}
}

func main() {
	var settingsPath *fayl.Path = gopolutils.Must(setup.Configuration(settings.Path))

	var configuration settings.Settings = settings.Read(settingsPath)

	var enviornmentFile *fayl.Path = gopolutils.Must(setup.Configuration(fayl.PathFrom(configuration.EnviornmentFilename)))
	var enviornmentVariables enviornment.Enviornment = enviornment.From(enviornmentFile)

	var token string = enviornment.Get(enviornmentVariables, configuration.Twitch.TokenKey)

	setToday(configuration.Twitch.Today.Topic)
	setLinks(configuration.Socials)

	var ircUrl url.URL = urlParse(fmt.Sprintf("%s://%s:%d", configuration.Twitch.Scheme, configuration.Twitch.BaseUrl, configuration.Twitch.Port))
	var authMessage string = fmt.Sprintf("%s %s", models.Authenticate, token)
	var nameMessage string = fmt.Sprintf("%s %s", models.Name, configuration.BotName)
	var joinMessage string = fmt.Sprintf("%s #%s", models.Join, configuration.Twitch.ChannelName)

	var logger *gopolutils.Logger = gopolutils.NewLogger(configuration.BotName, gopolutils.Info)
	var logFile *fayl.Path = gopolutils.Must(setup.Logging(configuration.Logging.Folder, configuration.Logging.Format))
	var except *gopolutils.Exception = logger.FileOnly(logFile.String())
	if except != nil {
		panic(except)
	}

	var irc *irc.IRC = irc.New(ircUrl.String())
	defer irc.Close()

	var bot *bot.Bot = bot.New(irc, logger, configuration)
	writeMessages(bot, authMessage, nameMessage, joinMessage)

	except = command.Setup()
	if except != nil {
		panic(except)
	}

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
