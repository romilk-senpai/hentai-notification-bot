package cmd

import (
	"flag"
	tgclient "hentai-notification-bot-re/client/telegram"
	tgcontroller "hentai-notification-bot-re/controller/telegram"
	parseflag "hentai-notification-bot-re/flag"
	"hentai-notification-bot-re/handler"
	"hentai-notification-bot-re/lib/e/config"
	"hentai-notification-bot-re/parser"
	"hentai-notification-bot-re/parser/hentaifox"
	"hentai-notification-bot-re/parser/nhentai"
	jsonrepository "hentai-notification-bot-re/repository/json"
	"log"
)

const (
	tgBotHost     = "api.telegram.org"
	nhentaiHost   = "nhentai.net"
	hentaifoxHost = "hentaifox.com"
	batchSize     = 100
)

func Main() {
	flag.Parse()

	if parseflag.TgToken == "" {
		log.Fatal("Telegram Token is empty")
	}

	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	parsers := []parser.Parser{nhentai.New(nhentaiHost, cfg), hentaifox.New(hentaifoxHost)}

	var repo jsonrepository.JsonRepository[*tgcontroller.UserInfo]
	tgRepo := repo.New("local-cache")
	tgClient := tgclient.New(tgBotHost, parseflag.TgToken)
	tgController := tgcontroller.New(
		tgClient,
		tgRepo,
		parsers,
	)

	var h handler.Handler

	if parseflag.WithWebhook {
		h = handler.NewTgHanlder(cfg, tgClient, tgController)
	} else {
		tgClient.SetWebhook("")
		h = handler.NewLocalHandler(tgController, tgController, batchSize)
	}

	log.Println("service is started")

	if err := h.Run(); err != nil {
		log.Fatal("service is stopped ", err)
	}
}
