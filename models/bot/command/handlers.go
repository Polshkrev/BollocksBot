package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/Polshkrev/BollocksBot/models/enviornment"
	"github.com/Polshkrev/gopolutils"
)

// Concurrently obtain the value stored at the given key within the system's enviornment.
func getConcurrent(key string, resultChannel chan<- string, errorChannel chan<- *gopolutils.Exception) {
	defer close(resultChannel)
	defer close(errorChannel)
	if !enviornment.Check(key) {
		errorChannel <- gopolutils.NewNamedException(gopolutils.KeyError, "Can not find key \"%s\" in enviornment.", key)
		return
	}
	resultChannel <- os.Getenv(key)
	errorChannel <- nil
}

// Obtain the value stored at the given key within the system's enviornment.
// Returns the value stored at the given key within the system's enviornment.
func get(key string) string {
	var resultChannel chan string = make(chan string, 1)
	var errorChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception, 1)
	go getConcurrent(key, resultChannel, errorChannel)
	var except *gopolutils.Exception = <-errorChannel
	if except != nil {
		panic(except)
	}
	var result string = <-resultChannel
	return result
}

// Callback for the `ban` command.
func handleBan(argument string) string {
	return fmt.Sprintf("%s has been banned.", argument)
}

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
	return fmt.Sprintf("@%s %s", argument, get(strings.ToUpper(Today.String())))
}

// Callback for the `vod` command.
func handleVod(argument string) string {
	var link string = get(strings.ToUpper(Youtube.String()))
	return fmt.Sprintf("@%s %s", argument, link)
}

// Callback for an unknown command.
func handleUnknown(argument string) string {
	return fmt.Sprintf("Unknwown Command: \"%s\"", argument)
}
