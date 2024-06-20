package tgcontroller

import (
	"fmt"
	events "hentai-notification-bot-re/controller"
	"hentai-notification-bot-re/lib/e"
)

const (
	Start   = "start"
	TestAdd = "addTagGroup"
)

func (c *Controller) processCmd(event events.Event) error {

	command := event.CommandInfo.Command

	switch command {
	case Start:
		{
			return c.start(event)
		}
	case TestAdd:
		{
			return c.addTagGroupCmd(event)
		}
	default:
		return c.unknown(event)
	}
}

func (c *Controller) start(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return err
	}

	if !c.repository.Exists(event.UserHash) {
		newUserTags := NewTagMap()

		expr, err := processAddExpression(event.CommandInfo.Arguments)

		if err != nil {
			_ = c.client.SendMessage(meta.Update.Message.Chat.ID, fmt.Sprintf("Expression error; arg=%s", expr))

			return err
		}

		_, err = c.repository.Create(NewUserInfo(event.UserHash, meta.Update.Message.From.Username, meta.Update.Message.Chat.ID, newUserTags))

		if err != nil {
			return err
		}
	}

	return c.client.SendStandardMarkup(meta.Update.Message.Chat.ID)
}

func (c *Controller) addTagGroupCmd(event events.Event) (err error) {
	defer func() { err = e.WrapIfErr("can't process request", err) }()

	if err != nil {
		return err
	}

	return c.addTagGroup(event.UserHash, event.CommandInfo.Arguments)
}

func (c *Controller) unknown(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return err
	}

	return c.client.SendMessage(meta.Update.Message.Chat.ID, "Unknown command")
}

func processAddExpression(expression string) (string, error) {
	return expression, nil
}
