package command

import "github.com/Polshkrev/gopolutils"

// Name of a specific command.
type Name gopolutils.StringEnum

const (
	Available Name = "commands" // List all available commands.
	Hello     Name = "hello"    // Say hello.
	Ping      Name = "ping"     // Ping the bot.
	Today     Name = "today"    // Find out what the stream topic is today.
)
