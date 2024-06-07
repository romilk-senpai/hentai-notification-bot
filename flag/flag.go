package parseflag

import "flag"

var (
	TgToken string
)

func init() {
	flag.StringVar(&TgToken, "tg-token", "", "telegram bot token")
}
