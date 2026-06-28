package bot

/*
#cgo CFLAGS: -O3 -march=native

#define LEXER_IMPLEMENTATION
#include "../../lexer/lexer.h"

#include <stdlib.h> // free
*/
import "C"
import (
	"bufio"
	"fmt"
	"strings"
	"unsafe"

	"github.com/Polshkrev/BollocksBot/models"
	"github.com/Polshkrev/BollocksBot/models/bot/command"
	"github.com/Polshkrev/BollocksBot/models/irc"
	"github.com/Polshkrev/BollocksBot/settings"
	"github.com/Polshkrev/gopolutils"
)

const (
	crlf        string = "\r\n"                // Line feed charactor.
	pingMessage string = "PONG :tmi.twitch.tv" // Standard response to a health check.
	pingStamp   string = "[PING]"              // Stamp to log on a ping message.
)

// IRC Bot.
type Bot struct {
	client    *irc.IRC
	logger    *gopolutils.Logger
	setttings settings.Settings
}

// Contruct a new [Bot] with a given [IRC] and [settings.Settings].
// Returns a new [Bot] with a given [IRC] and [settings.Settings].
func New(client *irc.IRC, logger *gopolutils.Logger, settings settings.Settings) *Bot {
	var bot *Bot = new(Bot)
	bot.client = client
	bot.logger = logger
	bot.setttings = settings
	return bot
}

// Send an authentication message to the IRC client.
// If the message can not be written, the function panics with a [gopolutils.RuntimeError].
func (bot *Bot) Write(message string) {
	var except *gopolutils.Exception = bot.client.Write(message)
	if except != nil {
		panic(except)
	}
}

// Send a properly formatted IRC repsonse.
func (bot *Bot) Respond(message string) {
	bot.Write(fmt.Sprintf("%s #%s :%s", models.Message, bot.setttings.Twitch.ChannelName, message))
	bot.logger.Log(fmt.Sprintf("%s: %s", bot.setttings.BotName, message), gopolutils.Info)
}

func isEmpty(value string) bool {
	return value == "" || len(value) == 0
}

func commandToString(rawCommand string, rawArgument string, recipient string) string {
	if isEmpty(rawArgument) {
		return command.HandleCommand(rawCommand, recipient)
	}
	return command.HandleCommand(rawCommand, rawArgument)
}

// Handle each message that comes through the bot's IRC client.
// Returns true if the message can be parsed, else false.
func (bot *Bot) HandleMessage(message string) bool {
	var lexer C.lexer_t = C.lexer_init()

	var messageCString *C.char = C.CString(message)
	defer C.free(unsafe.Pointer(messageCString))

	C.lexer_set_source(&lexer, messageCString)

	var tokens C.token_array_t
	C.tokenize(&lexer, &tokens)

	var messageType C.message_t = C.message_init()

	if !C.parse_message(&tokens, &messageType) {
		return false
	}

	var text SizedString = SizedStringFrom(messageType.text)

	if text.IsEmpty() {
		sendPing(bot, &messageType)
		return true
	}

	var name string = SizedStringFrom(messageType.name).String()

	bot.logger.Log(fmt.Sprintf("%s: %s", name, text), gopolutils.Info)

	if !C.parse_command(&messageType) {
		return true
	}

	var rawCommand string = SizedStringFrom(messageType.command).Trim()
	var rawArgument string = SizedStringFrom(messageType.arguments).String()

	bot.Respond(commandToString(rawCommand, rawArgument, name))
	return true
}

// Trim the trailing line feed charactor.
// Returns the given message without the trailing line feed charactor.
func trimNewline(message string) string {
	return strings.TrimSuffix(message, crlf)
}

// Read the current IRC client's buffer.
// If the message can not be read, the function panics with a [gopolutils.RuntimeError].
func (bot *Bot) Read() {
	for {
		var readMessage string = trimNewline(gopolutils.Must(bot.client.Read()))
		if readMessage == "" || len(readMessage) == 0 {
			continue
		} else if !bot.HandleMessage(readMessage) {
			panic(gopolutils.NewNamedException(gopolutils.ValueError, "Can not read message: \"%s\".", readMessage))
		}
	}
}

// Wait for the input from the given buffer, then write that ouput to the IRC client.
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

// Send the healthcheck response through the bot's IRC client.
func sendPing(bot *Bot, message *C.message_t) {
	if SizedStringFrom(message.keyword).String() != string(models.Ping) {
		return
	}
	bot.Write(handlePing(bot.logger))
}

// Obtain the response to a health check.
// Returns pong.
func handlePing(logger *gopolutils.Logger) string {
	logger.Log(pingStamp, gopolutils.Info)
	return pingMessage
}
