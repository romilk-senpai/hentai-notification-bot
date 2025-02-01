package handler

import (
	events "hentai-notification-bot-re/controller"
	"log"
	"time"
)

type LocalHandler struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewLocalHandler(fetcher events.Fetcher, processor events.Processor, batchSize int) *LocalHandler {
	return &LocalHandler{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (h *LocalHandler) Run() error {
	for {
		gotEvents, err := h.fetcher.Fetch(h.batchSize)

		if err != nil {
			log.Printf("[ERR] handler: %s", err.Error())

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

func (h *LocalHandler) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %v", event.Type)

		if err := h.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
