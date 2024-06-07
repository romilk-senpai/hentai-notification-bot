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

	/*for _, mParser := range c.parsers {
		mangoes, err := mParser.ParseAll(event.Text)

		if err != nil {
			return e.Wrap("parser error", err)
		}

		var responseBuilder strings.Builder

		if len(mangoes) == 0 {
			responseBuilder.WriteString("no manga with the given tags found")
		}

		for _, manga := range mangoes {
			responseBuilder.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a>\n", manga.Url, manga.Name))
			//responseBuilder.WriteString(fmt.Sprintf("%s: %s", manga.Name, manga.Url))
		}

		if err := c.client.SendMessage(meta.ChatID, responseBuilder.String()); err != nil {
			return e.Wrap("can't send message", err)
		}
	}*/

	return nil
}
