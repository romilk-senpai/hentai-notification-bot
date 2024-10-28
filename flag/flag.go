package parseflag

import "flag"

var (
	TgToken     string
	ConfigPath  string
	WithWebhook bool
)

func init() {
	flag.StringVar(&TgToken, "tg-token", "", "telegram bot token")
	flag.StringVar(&ConfigPath, "config", "", "server config")
	flag.BoolVar(&WithWebhook, "with-webhook", false, "enable webhook mode")
}
