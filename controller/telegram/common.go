package tgcontroller

import (
	"errors"
	"fmt"
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

	_, err = c.repository.Update(userHash, userInfo)

	if err != nil {
		return err
	}

	return c.client.SendMessage(userInfo.ChatID, fmt.Sprintf("Successfully added %s!\nKeep doing what u're doing 😭", expr))
}
