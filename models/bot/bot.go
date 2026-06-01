package bot

import (
	"bufio"
	"fmt"

	"github.com/Polshkrev/BollocksBot/models/irc"
	"github.com/Polshkrev/BollocksBot/settings"
	"github.com/Polshkrev/gopolutils"
)

// IRC Bot.
type Bot struct {
	client    *irc.IRC
	setttings settings.Settings
}

// Contruct a new [Bot] with a given [IRC] and [settings.Settings].
// Returns a new [Bot] with a given [IRC] and [settings.Settings].
func New(client *irc.IRC, settings settings.Settings) *Bot {
	var bot *Bot = new(Bot)
	bot.client = client
	bot.setttings = settings
	return bot
}

// Send an authentication message to the irc client.
// If the message can not be written, the function panics with a [gopolutils.RuntimeError].
func (bot *Bot) Authenticate(message string) {
	var except *gopolutils.Exception = bot.client.Write(message)
	if except != nil {
		panic(except)
	}
}

// Send a login message to the irc client.
// If the message can not be written, the function panics with a [gopolutils.RuntimeError].
func (bot *Bot) Login(message string) {
	var except *gopolutils.Exception = bot.client.Write(message)
	if except != nil {
		panic(except)
	}
}

// Send a join message to the irc client.
// If the message can not be written, the function panics with a [gopolutils.RuntimeError].
func (bot *Bot) Join(message string) {
	var except *gopolutils.Exception = bot.client.Write(message)
	if except != nil {
		panic(except)
	}
}

// Read the current irc client's buffer.
// If the message can not be read, the function panics with a [gopolutils.RuntimeError].
func (bot *Bot) Read() {
	for {
		fmt.Print(gopolutils.Must(bot.client.Read()))
	}
}

// Wait for the input from the given buffer, then write that ouput to the irc client.
// If the buffer can not be read, the function panics with an [gopolutils.IOError].
// If the message can not be written, the function panics with a [gopolutils.RuntimeError].
func (bot *Bot) WaitForInput(reader *bufio.Reader) {
	for {
		var command string
		var readerError error
		command, readerError = reader.ReadString('\n')
		if readerError != nil {
			panic(gopolutils.NewNamedException(gopolutils.IOError, "%s\n", readerError.Error()))
		}
		var commandError *gopolutils.Exception = bot.client.Write(command)
		if commandError != nil {
			panic(commandError)
		}
	}

}
