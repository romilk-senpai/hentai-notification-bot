package tgcontroller

import (
	"fmt"
	events "hentai-notification-bot-re/controller"
	"hentai-notification-bot-re/lib/e"
)

const (
	Start   = "start"
	TestAdd = "testAdd"
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
			return c.testAdd(event)
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

	return c.client.SendStandardMarkup(meta.ChatID)
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
