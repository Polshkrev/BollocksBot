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
func handleBan(argument, sender string) string {
	return fmt.Sprintf("@%s \"%s\" has been banned.", sender, argument)
}

// Callback for the 'hello' command.
func handleHello(_, sender string) string {
	return fmt.Sprintf("Hello @%s", sender)
}

// Callback for the 'ping' command.
func handlePing(_, sender string) string {
	return fmt.Sprintf("@%s pong", sender)
}

// Callback for the 'today' command.
func handleToday(_, sender string) string {
	return fmt.Sprintf("@%s %s", sender, get(strings.ToUpper(Today.String())))
}

// Callback for the 'project' command.
func handleProject(_, sender string) string {
	var upperKey string = strings.ToUpper(Project.String())
	if !enviornment.Check(upperKey) {
		return fmt.Sprintf("@%s This is no project set today. Try the \"!today\" command.", sender)
	}
	return fmt.Sprintf("@%s %s", sender, get(upperKey))
}

// Callback for the `vod` command.
func handleVod(_, sender string) string {
	var link string = get(strings.ToUpper(Youtube.String()))
	return fmt.Sprintf("@%s %s", sender, link)
}

// Callback for an unknown command.
func handleUnknown(argument, sender string) string {
	return fmt.Sprintf("@%s Unknwown Command: \"%s\"", sender, argument)
}
