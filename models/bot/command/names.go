package command

import "github.com/Polshkrev/gopolutils"

// Name of a specific command.
type Name gopolutils.StringEnum

const (
	Ban       Name = "ban"      // A ban command.
	Available Name = "commands" // List all available commands.
	Hello     Name = "hello"    // Say hello.
	Ping      Name = "ping"     // Ping the bot.
	Today     Name = "today"    // Find out what the stream topic is today.
	Youtube   Name = "vod"      // Display the vod channel.
)

// Returns a string representation of a [Name].
func (name Name) String() string {
	return string(name)
}
