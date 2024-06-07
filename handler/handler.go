package handler

import (
	events "hentai-notification-bot-re/controller"
	"log"
	"time"
)

type Handler struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Handler {
	return Handler{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (h *Handler) Run() error {
	for {
		gotEvents, err := h.fetcher.Fetch(h.batchSize)

		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := h.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (h *Handler) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := h.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
