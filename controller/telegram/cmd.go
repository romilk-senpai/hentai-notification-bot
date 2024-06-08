package tgcontroller

import (
	"fmt"
	events "hentai-notification-bot-re/controller"
	"hentai-notification-bot-re/lib/e"
	"strings"
)

const (
	Test           = "test"
	TestAdd        = "testAdd"
	TestGetUpdates = "testGetUpdates"
)

func (c *Controller) processCmd(event events.Event) error {

	command := event.CommandInfo.Command

	switch command {
	case Test:
		{
			return c.test(event)
		}
	case TestAdd:
		{
			return c.testAdd(event)
		}
	case TestGetUpdates:
		{
			return c.getUpdates(event)
		}
	default:
		return c.unknown(event)
	}
}

func (c *Controller) test(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return e.Wrap(fmt.Sprintf("can't process command %s", event.CommandInfo.Command), err)
	}

	return c.client.SendMessage(meta.ChatID, fmt.Sprintf("test command; arg=%s", event.CommandInfo.Arguments))
}

func (c *Controller) testAdd(event events.Event) (err error) {
	meta, err := meta(event)

	defer func() { err = e.WrapIfErr("can't process request", err) }()

	if err != nil {
		return err
	}

	expr, err := processAddExpression(event.CommandInfo.Arguments)

	if err != nil {
		_ = c.client.SendMessage(meta.ChatID, fmt.Sprintf("Expression error; arg=%s", expr))

		return err
	}

	if c.repository.Exists(event.UserHash) {
		userInfo, err := c.repository.Read(event.UserHash)

		if err != nil {
			return err
		}

		if userInfo.SubscribedTags == nil {
			userInfo.SubscribedTags = make(map[string]map[string]int)
			userInfo.SubscribedTags[expr] = make(map[string]int)
		}

		_, err = c.repository.Update(event.UserHash, userInfo)
	} else {
		newUserTags := make(map[string]map[string]int)
		newUserTags[expr] = make(map[string]int)

		expr, err := processAddExpression(event.CommandInfo.Arguments)

		if err != nil {
			_ = c.client.SendMessage(meta.ChatID, fmt.Sprintf("Expression error; arg=%s", expr))

			return err
		}

		_, err = c.repository.Create(UserInfo{
			Uuid:           event.UserHash,
			Username:       meta.Username,
			ChatID:         meta.ChatID,
			SubscribedTags: newUserTags,
		})

		if err != nil {
			return err
		}
	}

	return c.client.SendMessage(meta.ChatID, fmt.Sprintf("Added; arg=%s", expr))
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

func (c *Controller) unknown(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return err
	}

	return c.client.SendMessage(meta.ChatID, "Unknown command")
}

func processAddExpression(expression string) (string, error) {
	return expression, nil
}
