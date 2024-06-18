package tgcontroller

import (
	"errors"
	tgclient "hentai-notification-bot-re/client/telegram"
	events "hentai-notification-bot-re/controller"
)

func (c *Controller) processCallback(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return err
	}

	data := meta.Update.CallbackQuery.ParseCallbackData()

	user, err := c.repository.Read(event.UserHash)

	switch data.Key {
	case "addTagGroup":
		{
			return c.addTagGroup(data, user)
		}
	case "deleteTags":
		{
			return c.deleteTagGroup(meta.Update.CallbackQuery, data, user)
		}

	case "cancelManage":
		{
			return c.cancelProcessTags(meta.Update.CallbackQuery)
		}
	}

	return nil
}

func (c *Controller) addTagGroup(data *tgclient.CallbackData, user UserInfo) error {
	return nil
}

func (c *Controller) cancelProcessTags(query *tgclient.CallbackQuery) error {
	return c.client.DeleteMessage(query.Message.Chat.ID, query.Message.ID)
}

func (c *Controller) deleteTagGroup(query *tgclient.CallbackQuery, data *tgclient.CallbackData, userInfo UserInfo) error {
	exists := userInfo.SubscribedTags.SubscribedToTag(data.Value)

	if !exists {
		return errors.New("trying to delete non-existent tag group")
	}

	userInfo.SubscribedTags.Delete(data.Value)

	_, err := c.repository.Update(userInfo.Uuid, userInfo)

	if err != nil {
		return err
	}

	return c.client.EditTagManager(query.Message.Chat.ID, query.Message.ID, userInfo.SubscribedTags.Tags)
}
