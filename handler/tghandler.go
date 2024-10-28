package handler

import (
	tgclient "hentai-notification-bot-re/client/telegram"
	events "hentai-notification-bot-re/controller"
	tgcontroller "hentai-notification-bot-re/controller/telegram"
	"hentai-notification-bot-re/lib/e/config"
	"log"
	"net/http"
)

type TgHandler struct {
	tgClient  *tgclient.Client
	processor events.Processor
}

func NewTgHanlder(tgClient *tgclient.Client, processor events.Processor) TgHandler {
	return TgHandler{
		tgClient:  tgClient,
		processor: processor,
	}
}

func (h *TgHandler) Run() error {
	cfg, err := config.Load()

	if err != nil {
		return err
	}

	err = h.tgClient.SetWebhook(cfg.HTTPServer.Host + "/tgevent")

	if err != nil {
		log.Fatal("failed to set webhook ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/tgevent", tgcontroller.NewHandler(h.processor))

	srv := &http.Server{
		Handler:      mux,
		Addr:         cfg.HTTPServer.Address,
		WriteTimeout: cfg.HTTPServer.Timeout,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
