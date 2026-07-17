package command

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
)

var (
	commandEntries collections.Mapping[string, Caller] = collections.NewMap[string, Caller]() // Command lookup table.
)

type Caller func(argument string, sender string) string // Callback type alias.

// Setup the command lookup table.
// If the commands can not be setup, a [gopolutils.KeyError] is returned.
func Setup() *gopolutils.Exception {
	var except *gopolutils.Exception = commandEntries.Insert(Available.String(), AvailableCommands)
	if except != nil {
		return except
	}
	except = commandEntries.Insert(Ping.String(), handlePing)
	if except != nil {
		return except
	}
	except = commandEntries.Insert(Hello.String(), handleHello)
	if except != nil {
		return except
	}
	except = commandEntries.Insert(Today.String(), handleToday)
	if except != nil {
		return except
	}
	except = commandEntries.Insert(Project.String(), handleProject)
	if except != nil {
		return except
	}
	except = commandEntries.Insert(Youtube.String(), handleVod)
	if except != nil {
		return except
	}
	except = commandEntries.Insert(Ban.String(), handleBan)
	if except != nil {
		return except
	}
	return nil
}

// Call a command's callback of a given name with a given argument.
// Returns a new message to write to the chat.
func HandleCommand(name, argument, sender string) string {
	if !commandEntries.HasKey(name) {
		return handleUnknown(fmt.Sprintf("!%s", name), sender)
	}
	return (*gopolutils.Must(commandEntries.At(name)))(argument, sender)
}

// Sort the given commands alphabetically.
// The given command slice is modified.
func sortCommands(names *[]string) {
	slices.Sort(*names)
}

// Prepend a bang charactor to each name in the given name slice.
// Returns a new slice based on the prepended charactor.
func formatCommand(names []string) []string {
	var result []string
	var i int
	for i = range names {
		result = append(result, fmt.Sprintf("!%s", names[i]))
	}
	return result
}

// Represent all available command options as a string.
// Returns a string representation of all available command options.
func AvailableCommands(_, _ string) string {
	var keys []string = commandEntries.Keys()
	sortCommands(&keys)
	return strings.Join(formatCommand(keys), "; ")
}
