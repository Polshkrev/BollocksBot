package enviornment

import (
	"os"
	"strings"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/collections"
	"github.com/Polshkrev/gopolutils/fayl"
	"github.com/joho/godotenv"
)

type Enviornment collections.Mapping[string, string] // Representation of the system envionment.

const (
	Seperator string = "=" // Seperator between the system variables.
)

type Variable collections.Pair[string, string] // Representation of a single system variable.

func check(key string, resultChannel chan<- bool) {
	defer close(resultChannel)
	var ok bool
	_, ok = os.LookupEnv(key)
	resultChannel <- ok
}

// Check if the given key is stored in the system's enviornment.
// Returns true if the given key is set within the system's enviornment.
func Check(key string) bool {
	var resultChannel chan bool = make(chan bool, 1)
	go check(key, resultChannel)
	var result bool = <-resultChannel
	return result
}

// Concurrently set a given value at a given key within the system's enviornment.
func set(key, value string, errorChannel chan<- *gopolutils.Exception) {
	defer close(errorChannel)
	var upperValue string = strings.ToUpper(key)
	if Check(upperValue) {
		errorChannel <- nil
		return
	}
	var setError error = os.Setenv(upperValue, value)
	if setError != nil {
		errorChannel <- gopolutils.NewNamedException(gopolutils.KeyError, "%s", setError)
		return
	}
	errorChannel <- nil
}

// Set the given value at the given key within the system enviornment.
// If the given key already exists within the system enviornment, the function returns with a nil exception value.
func Set(key, value string) *gopolutils.Exception {
	var errorChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception, 1)
	go set(key, value, errorChannel)
	var except *gopolutils.Exception = <-errorChannel
	return except
}

// Concurrently obtain the given key from the given enviornment.
func get(envionment Enviornment, key string, resultChannel chan<- string, errorChannel chan<- *gopolutils.Exception) {
	defer close(resultChannel)
	defer close(errorChannel)
	if Check(key) {
		resultChannel <- os.Getenv(key)
		errorChannel <- nil
		return
	}
	var result *string
	var except *gopolutils.Exception
	result, except = envionment.At(key)
	resultChannel <- *result
	errorChannel <- except
}

// Obtain the value set at the given key within the system variables.
// Returns the value set at the given key within the system variables.
func Get(enviornment Enviornment, key string) string {
	var resultChannel chan string = make(chan string, 1)
	var errorChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception, 1)
	go get(enviornment, key, resultChannel, errorChannel)
	var except *gopolutils.Exception = <-errorChannel
	if except != nil {
		panic(except)
	}
	var result string = <-resultChannel
	return result
}

// Load a the system enviornment from a given file.
// Returns an [Enviornment] from the given [fayl.Path].
// If the given file does not exist, the enviornment is loaded from the system's internal enviornment.
// If the .env file can not be read, this function `panics` with an [gopolutils.OSError].
// If the key is already in the mapping, instead of just quietly not inserting into the mapping, this function `panics` with a [gopolutils.KeyError].
func From(file *fayl.Path) Enviornment {
	if !file.Exists() {
		return Load()
	}
	var result Enviornment = collections.NewMap[string, string]()
	var raw map[string]string
	var readError error
	raw, readError = godotenv.Read(file.String())
	if readError != nil {
		panic(gopolutils.NewNamedException(gopolutils.OSError, "%s\n", readError.Error()))
	}
	var key, value string
	for key, value = range raw {
		var insertExcept *gopolutils.Exception = result.Insert(key, value)
		if insertExcept != nil {
			panic(insertExcept)
		}
	}
	return result
}

// Concurrently split the given raw enviornment variable based on the default seperator.
func split(variable string, resultChannel chan<- Variable, errorChannel chan<- *gopolutils.Exception) {
	defer close(resultChannel)
	defer close(errorChannel)
	var split []string = strings.Split(variable, Seperator)
	if len(split) > 2 {
		errorChannel <- gopolutils.NewNamedException(gopolutils.ValueError, "Can not seperate the variable \"%s\".", variable)
		return
	}
	resultChannel <- Variable(*collections.NewPair(split[0], split[1]))
	errorChannel <- nil
}

// Split a given string representation of a given enviornment variable.
// Returns the given string representation as a [Variable].
// If the variable can not be split, a [gopolutils.ValueError] is returned.
func Split(variable string) (Variable, *gopolutils.Exception) {
	var resultChannel chan Variable = make(chan Variable, 1)
	var errorChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception, 1)
	go split(variable, resultChannel, errorChannel)
	var except *gopolutils.Exception = <-errorChannel
	var result Variable = <-resultChannel
	return result, except
}

// Concurrently load the raw system enviornment.
func loadResult(resultChannel chan<- []string) {
	defer close(resultChannel)
	resultChannel <- os.Environ()
}

// Load the raw system enviornment.
// Returns a slice of raw enviornment variables.
func loadRaw() []string {
	var resultChannel chan []string = make(chan []string, 1)
	go loadResult(resultChannel)
	var result []string = <-resultChannel
	return result
}

// Concurrently load the system enviornment as an [Enviornment].
func assign(resultChannel chan<- Enviornment, errorChannel chan<- *gopolutils.Exception) {
	defer close(resultChannel)
	defer close(errorChannel)
	var raw []string = loadRaw()
	var result Enviornment = Enviornment(collections.NewMap[string, string]())
	var i int
	for i = range raw {
		var variable collections.Pair[string, string] = collections.Pair[string, string](gopolutils.Must(Split(raw[i])))
		var except *gopolutils.Exception = result.Insert(*variable.First(), *variable.Second())
		if except != nil {
			errorChannel <- except
			return
		}
	}
	resultChannel <- result
	errorChannel <- nil
}

// Load the system enviornment variables as an [Enviornment].
// Returns the system enviornment variables as an [Enviornment].
// If the enviornment can not be parsed or loaded, a [gopolutils.KeyError] is returns.
func load() (Enviornment, *gopolutils.Exception) {
	var resultChannel chan Enviornment = make(chan Enviornment, 1)
	var errorChannel chan *gopolutils.Exception = make(chan *gopolutils.Exception)
	go assign(resultChannel, errorChannel)
	var except *gopolutils.Exception = <-errorChannel
	if except != nil {
		return nil, except
	}
	var result Enviornment = <-resultChannel
	return result, nil
}

// Load the system enviornment from the internal enviornment.
// Returns an [Enviornment] based on the system's internal enviornment.
func Load() Enviornment {
	return gopolutils.Must(load())
}
