package command

import "fmt"

// Callback for the 'hello' command.
func handleHello(argument string) string {
	return fmt.Sprintf("Hello @%s", argument)
}

// Callback for the 'ping' command.
func handlePing(argument string) string {
	return fmt.Sprintf("@%s pong", argument)
}

// Callback for an unknown command.
func handleUnknown(argument string) string {
	return fmt.Sprintf("Unknwown Command: '%s'", argument)
}
