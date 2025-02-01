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
	cfg       *config.Config
}

func NewTgHanlder(cfg *config.Config, tgClient *tgclient.Client, processor events.Processor) *TgHandler {
	return &TgHandler{
		cfg:       cfg,
		tgClient:  tgClient,
		processor: processor,
	}
}

func (h *TgHandler) Run() error {
	err := h.tgClient.SetWebhook(h.cfg.HTTPServer.Host + "/tgevent")

	if err != nil {
		log.Fatal("failed to set webhook ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/tgevent", tgcontroller.NewHandler(h.processor))

	srv := &http.Server{
		Handler:      mux,
		Addr:         h.cfg.HTTPServer.Address,
		WriteTimeout: h.cfg.HTTPServer.Timeout,
		ReadTimeout:  h.cfg.HTTPServer.Timeout,
		IdleTimeout:  h.cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
