package irc

import (
	"github.com/Polshkrev/gopolutils"
	"github.com/gorilla/websocket"
)

// An IRC client.
type IRC struct {
	connection *websocket.Conn
}

// Contruct a new IRC client with a given url. This constructor calls the private dial method.
// Returns a new IRC client.
func New(url string) *IRC {
	var irc *IRC = new(IRC)
	irc.dial(url)
	return irc
}

// Dial a websocket with a given url.
func (irc *IRC) dial(url string) {
	var connection *websocket.Conn
	var dialError error
	connection, _, dialError = websocket.DefaultDialer.Dial(url, nil)
	if dialError != nil {
		panic(gopolutils.NewNamedException(gopolutils.RuntimeError, "%s\n", dialError.Error()))
	}
	irc.connection = connection
}

// Concurrently write a given message to a given websocket connection.
func writeConcurrent(connection *websocket.Conn, message []byte, errorChannel chan<- error) {
	defer close(errorChannel)
	var readError error = connection.WriteMessage(websocket.TextMessage, message)
	errorChannel <- readError
}

// Write a given message to the irc connection.
// If the message can not be written, a [gopolutils.RuntimeError] is returned.
func (irc *IRC) Write(message string) *gopolutils.Exception {
	var errorChannel chan error = make(chan error, 1)
	go writeConcurrent(irc.connection, []byte(message), errorChannel)
	var except error = <-errorChannel // Blocking operation
	if except != nil {
		return gopolutils.NewNamedException(gopolutils.RuntimeError, "%s\n", except.Error())
	}
	return nil
}

// Concurrently read from a given websocket connection.
func readConcurrent(connection *websocket.Conn, resultChannel chan<- []byte, errorChannel chan<- error) {
	defer close(resultChannel)
	defer close(errorChannel)
	var result []byte
	var readError error
	_, result, readError = connection.ReadMessage()
	resultChannel <- result
	errorChannel <- readError
}

// Read from the irc connection.
// Returns the string representation of the irc connection.
// If the connection can not be read, a [gopolutils.RuntimeError] is returned with an empty string.
func (irc *IRC) Read() (string, *gopolutils.Exception) {
	var resultChannel chan []byte = make(chan []byte, 1)
	var errorChannel chan error = make(chan error, 1)
	go readConcurrent(irc.connection, resultChannel, errorChannel)
	var result []byte = <-resultChannel
	var except error = <-errorChannel
	if except != nil {
		return "", gopolutils.NewNamedException(gopolutils.RuntimeError, "%s\n", except.Error())
	}
	return string(result), nil
}

// Close the irc connection.
// If closing the connection fails, a [gopolutils.RuntimeError] is returned.
func (irc *IRC) Close() *gopolutils.Exception {
	var closeError error = irc.connection.Close()
	if closeError != nil {
		return gopolutils.NewNamedException(gopolutils.RuntimeError, "%s\n", closeError.Error())
	}
	return nil
}
