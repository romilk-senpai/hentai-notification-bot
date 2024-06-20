package tgcontroller

import (
	"fmt"
	tgclient "hentai-notification-bot-re/client/telegram"
	events "hentai-notification-bot-re/controller"
	"hentai-notification-bot-re/lib/e"
	"log"
	"strings"
)

func (c *Controller) processMessage(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't process message", err) }()

	meta, err := meta(event)

	log.Printf("message from %s: %s", meta.Update.Message.From.Username, event.Text)

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
	default:
		{
			return c.processPlaneMessage(event)
		}
	}
}

func (c *Controller) processPlaneMessage(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't process plane message", err) }()

	meta, err := meta(event)

	if err != nil {
		return err
	}

	userInfo, err := c.repository.Read(event.UserHash)

	if err != nil {
		return err
	}

	if !userInfo.AddingTags {
		return nil
	}

	userInfo.AddingTags = false
	_, err = c.repository.Update(userInfo.Uuid, userInfo)

	if err != nil {
		return err
	}

	err = c.addTagGroup(event.UserHash, meta.Update.Message.Text)

	if err != nil {
		return err
	}

	userInfo, err = c.repository.Read(userInfo.Uuid)

	tagManagerMsgID, err := c.client.SendTagManager(userInfo.ChatID, userInfo.SubscribedTags.Tags)

	if err != nil {
		return err
	}

	userInfo.ManagerMessageID = tagManagerMsgID

	_, err = c.repository.Update(userInfo.Uuid, userInfo)

	return err
}

func (c *Controller) getUpdates(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	meta, err := meta(event)

	if err != nil {
		return err
	}

	userInfo, err := c.repository.Read(event.UserHash)

	if err != nil {
		return err
	}

	var responseBuilder strings.Builder

	userInfo.SubscribedTags.ForEach(func(tagGroup string, parserMap map[string]int) error {
		responseBuilder.WriteString(fmt.Sprintf("Updates on <b><u>%s</u></b>\n", tagGroup))
		for _, parser := range c.parsers {
			remoteQuantity, err := parser.ParseQuantity(tagGroup)

			if err != nil {
				return e.Wrap("parser error", err)
			}

			savedQuantity, ok := parserMap[parser.ParserName()]

			var parserResponseBuilder strings.Builder

			parserResponseBuilder.WriteString(fmt.Sprintf("\t<a href=\"%s\">%s (%d)</a>\n", parser.QueryToLink(tagGroup), parser.ParserName(), remoteQuantity-savedQuantity))

			if !ok {
				savedQuantity = remoteQuantity

				parserMap[parser.ParserName()] = savedQuantity

				if _, err = c.repository.Update(event.UserHash, userInfo); err != nil {
					return err
				}
			}

			if savedQuantity >= remoteQuantity {
				responseBuilder.WriteString(parserResponseBuilder.String() + "\n")

				continue
			}

			mangoes, err := parser.ParseAll(tagGroup)

			if err != nil {
				return err
			}

			if len(mangoes) == 0 {
				parserResponseBuilder.WriteString("\tNo manga with the given tags found\n")

				responseBuilder.WriteString(parserResponseBuilder.String() + "\n")

				continue
			}

			for i := 0; i < remoteQuantity-savedQuantity; i++ {
				manga := mangoes[i]

				parserResponseBuilder.WriteString(fmt.Sprintf("\t<a href=\"%s\">%s</a>\n", manga.Url, manga.Name))
			}

			responseBuilder.WriteString(parserResponseBuilder.String() + "\n")

			parserMap[parser.ParserName()] = remoteQuantity

			if _, err = c.repository.Update(event.UserHash, userInfo); err != nil {
				return err
			}
		}

		return nil
	})

	if err = c.client.SendMessage(meta.Update.Message.Chat.ID, responseBuilder.String()); err != nil {
		return err
	}

	return nil
}

func (c *Controller) manageTags(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't manage tags", err) }()

	if err != nil {
		return err
	}

	userInfo, err := c.repository.Read(event.UserHash)

	if err != nil {
		return err
	}

	tagManagerMsgID, err := c.client.SendTagManager(userInfo.ChatID, userInfo.SubscribedTags.Tags)

	if err != nil {
		return err
	}

	if userInfo.ManagerMessageID != -1 {
		err = c.client.DeleteMessage(userInfo.ChatID, userInfo.ManagerMessageID)

		if err != nil {
			return err
		}
	}

	userInfo.ManagerMessageID = tagManagerMsgID

	_, err = c.repository.Update(userInfo.Uuid, userInfo)

	return err
}
