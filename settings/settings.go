package settings

import (
	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/fayl"
)

const (
	Folder string      = "settings" // Parent folder of the settings file.
	Name   string      = "settings" // Name of the settings file.
	Suffix fayl.Suffix = fayl.Json  // [fayl.Suffix] of the settings file.
)

var (
	Path *fayl.Path = fayl.PathFromParts(Folder, Name, Suffix) // Full [fayl.Path] of the settings.
)

type TwitchSettings struct {
	TokenKey    string `json:"tokenKey"`    // Key to lookup the twitch token within the enviornment mapping.
	Scheme      string `json:"scheme"`      // Scheme of the twitch url.
	BaseUrl     string `json:"baseUrl"`     // Base url of the twitch irc client.
	Port        uint16 `json:"port"`        // Port of the twitch irc client.
	ChannelName string `json:"channelName"` // Target channel name.
}

// Global settings of the application.
type Settings struct {
	BotName             string         `json:"botName"`             // name of the bot.
	EnviornmentFilename string         `json:"enviornmentFilename"` // Name of the .env file.
	Twitch              TwitchSettings `json:"twitch"`              // Encapsulation of the twitch-specific settings.
}

// Read a [Settings] object from a given [fayl.Path].
func Read(file *fayl.Path) Settings {
	return *gopolutils.Must(fayl.ReadObject[Settings](file))
}
