package setup

import "github.com/Polshkrev/goserialize"

const (
	userFolder string = "polshkrev"   // Group configuration folder name.
	folder     string = "BollocksBot" // Configuration folder name. This should be set to the name of the bot.
)

var (
	// Default settings to copy to the configuration file.
	defaultSettings goserialize.Object = goserialize.Object{
		"botName":             "BollocksBot",
		"enviornmentFilename": ".env",
		"twitch": goserialize.Object{
			"tokenKey":    "TWITCH_TOKEN",
			"scheme":      "wss",
			"baseUrl":     "irc-ws.chat.twitch.tv",
			"port":        443,
			"channelName": "polshkrev",
			"today": goserialize.Object{
				"topic":   "Making the chat bot.",
				"project": "https://github.com/Polshkrev/BollocksBot",
			},
		},
		"socials": goserialize.Object{
			"vod": "https://www.youtube.com/@polshkrev",
		},
		"logging": goserialize.Object{
			"folder": "logs",
			"format": "2006-01",
		},
	}
)
