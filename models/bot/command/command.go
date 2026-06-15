package command

import (
	"fmt"
	"strings"

	"github.com/Polshkrev/gopolutils/collections"
)

var (
	commandEntries collections.Collection[Command] = collections.NewArray[Command]() // Command lookup table.
)

type CommandCaller func(argument string) string // Command callback type alias.

// Representation of a command in chat.
type Command struct {
	name   string
	caller CommandCaller
}

// Construct a new command with a given name and callback.
// Returns a new command with a given name and callback.
func New(name string, caller CommandCaller) Command {
	var command *Command = new(Command)
	command.name = name
	command.caller = caller
	return *command
}

// Setup the command lookup table.
func Setup() {
	commandEntries.Append(New(string(Available), AvailableCommands))
	commandEntries.Append(New(string(Ping), handlePing))
	commandEntries.Append(New(string(Hello), handleHello))
	commandEntries.Append(New(string(Today), handleToday))
}

// Call a command's callback of a given name with a given argument.
// Returns a new message to write to the chat.
func HandleCommand(name string, argument string) string {
	var command Command
	for _, command = range commandEntries.Collect() {
		if command.name != name {
			continue
		}
		return command.caller(argument)
	}
	return handleUnknown(name)
}

// Represent all available command options as a string.
// Returns a string representation of all available command options.
func AvailableCommands(_ string) string {
	var command Command
	var names []string
	for _, command = range commandEntries.Collect() {
		names = append(names, fmt.Sprintf("!%s", command.name))
	}
	return strings.Join(names, "; ")
}
