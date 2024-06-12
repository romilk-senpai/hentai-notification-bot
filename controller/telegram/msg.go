package tgcontroller

import (
	"fmt"
	tgclient "hentai-notification-bot-re/client/telegram"
	events "hentai-notification-bot-re/controller"
	"hentai-notification-bot-re/lib/e"
	"log"
	"strings"
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

	switch event.Text {
	case tgclient.GetMangaUpdatesMessage:
		{
			return c.getUpdates(event)
		}
	case tgclient.ManageTagsMessage:
		{
			return c.manageTags(event)
		}
	}

	return nil
}

func (c *Controller) getUpdates(event events.Event) (err error) {
	meta, err := meta(event)

	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	userInfo, err := c.repository.Read(event.UserHash)

	if err != nil {
		return err
	}

	for tagGroup, parserMap := range userInfo.SubscribedTags {
		for _, mParser := range c.parsers {
			remoteQuantity, err := mParser.ParseQuantity(tagGroup)

			if err != nil {
				return e.Wrap("parser error", err)
			}

			savedQuantity, ok := parserMap[mParser.ParserName()]

			var responseBuilder strings.Builder

			responseBuilder.WriteString(fmt.Sprintf("%s updates:\n", mParser.ParserName()))

			if !ok {
				savedQuantity = remoteQuantity

				parserMap[mParser.ParserName()] = savedQuantity

				if _, err = c.repository.Update(event.UserHash, userInfo); err != nil {
					return err
				}
			}

			if savedQuantity >= remoteQuantity {
				responseBuilder.WriteString("no updates")

				if err = c.client.SendMessage(meta.ChatID, responseBuilder.String()); err != nil {
					return err
				}

				continue
			}

			mangoes, err := mParser.ParseAll(tagGroup)

			if err != nil {
				return err
			}

			if len(mangoes) == 0 {
				responseBuilder.WriteString("no manga with the given tags found")

				_ = c.client.SendMessage(meta.ChatID, responseBuilder.String())

				continue
			}

			for i := 0; i < remoteQuantity-savedQuantity; i++ {
				manga := mangoes[i]

				responseBuilder.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a>\n", manga.Url, manga.Name))
			}

			if err = c.client.SendMessage(meta.ChatID, responseBuilder.String()); err != nil {
				return err
			}

			parserMap[mParser.ParserName()] = remoteQuantity

			if _, err = c.repository.Update(event.UserHash, userInfo); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Controller) manageTags(event events.Event) (err error) {
	meta, err := meta(event)

	defer func() { err = e.WrapIfErr("can't manage tags", err) }()

	if err != nil {
		return err
	}

	userInfo, err := c.repository.Read(event.UserHash)

	if err != nil {
		return err
	}

	tagGroups := make([]string, 0, len(userInfo.SubscribedTags))

	for key := range userInfo.SubscribedTags {
		tagGroups = append(tagGroups, key)
	}

	return c.client.SendTagManager(meta.ChatID, tagGroups)
}
