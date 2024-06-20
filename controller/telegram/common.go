package tgcontroller

import (
	"errors"
	"fmt"
	"hentai-notification-bot-re/lib/e"
)

func (c *Controller) addTagGroup(userHash string, tagGroup string) error {
	if !c.repository.Exists(userHash) {
		return errors.New("user not found")
	}

	userInfo, err := c.repository.Read(userHash)

	if err != nil {
		return err
	}

	expr, err := processAddExpression(tagGroup)

	if err != nil {
		_ = c.client.SendMessage(userInfo.ChatID, fmt.Sprintf("Expression error; arg=%s", expr))

		return err
	}

	if userInfo.SubscribedTags == nil {
		userInfo.SubscribedTags = NewTagMap()
	}

	_, exists := userInfo.SubscribedTags.Get(expr)

	if !exists {
		userInfo.SubscribedTags.Set(expr, make(map[string]int))
	}

	for _, parser := range c.parsers {
		remoteQuantity, err := parser.ParseQuantity(tagGroup)

		if err != nil {
			return e.Wrap("parser error", err)
		}

		parserMap, _ := userInfo.SubscribedTags.Get(tagGroup)

		parserMap[parser.ParserName()] = remoteQuantity
	}

	_, err = c.repository.Update(userHash, userInfo)

	if err != nil {
		return err
	}

	return c.client.SendMessage(userInfo.ChatID, fmt.Sprintf("Successfully added <b>%s</b>!\nKeep doing what u're doing 😭", expr))
}
