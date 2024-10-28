package parseflag

import "flag"

var (
	TgToken    string
	ConfigPath string
)

func init() {
	flag.StringVar(&TgToken, "tg-token", "", "telegram bot token")
	flag.StringVar(&ConfigPath, "config", "", "Server config")
}
