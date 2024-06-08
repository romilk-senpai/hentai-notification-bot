package tgcontroller

import (
	events "hentai-notification-bot-re/controller"
	"hentai-notification-bot-re/lib/e"
	"log"
)

func (c *Controller) processMessage(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return e.Wrap("can't process message", err)
	}

	log.Printf("message from %s: %s", meta.Username, event.Text)

	if len(event.Text) == 0 {
		return nil
	}

	return nil
}
