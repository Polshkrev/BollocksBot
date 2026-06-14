package command

import (
	"fmt"
	"os"
)

// Callback for the 'hello' command.
func handleHello(argument string) string {
	return fmt.Sprintf("Hello @%s", argument)
}

// Callback for the 'ping' command.
func handlePing(argument string) string {
	return fmt.Sprintf("@%s pong", argument)
}

// Callback for the 'today' command
func handleToday(argument string) string {
	return fmt.Sprintf("@%s %s", argument, os.Getenv("TODAY_COMMAND"))
}

// Callback for an unknown command.
func handleUnknown(argument string) string {
	return fmt.Sprintf("Unknwown Command: '%s'", argument)
}
