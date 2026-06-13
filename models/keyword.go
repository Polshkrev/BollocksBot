package models

import "github.com/Polshkrev/gopolutils"

// Representation of a keyword to be sent through the IRC client.
type KeyWord gopolutils.StringEnum

const (
	Authenticate KeyWord = "PASS"    // Send the authentication token.
	Name         KeyWord = "NICK"    // Send the name of the IRC client.
	Join         KeyWord = "JOIN"    // Join a specific channel.
	Leave        KeyWord = "PART"    // Leave a specific channel.
	Message      KeyWord = "PRIVMSG" // Send a loggable message to through the IRC client.
	Request      KeyWord = "CAPREQ"  // Request the capabilities of the IRC client.
	Ping         KeyWord = "PING"    // Health check sent from twitch.
)
