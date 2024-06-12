package tgcontroller

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	tgclient "hentai-notification-bot-re/client/telegram"
	events "hentai-notification-bot-re/controller"
	"hentai-notification-bot-re/lib/e"
	"hentai-notification-bot-re/parser"
	"hentai-notification-bot-re/repository"
	"log"
)

type Controller struct {
	client     *tgclient.Client
	offset     int
	repository repository.Repository[UserInfo]
	parsers    []parser.Parser
}

func New(client *tgclient.Client, repository repository.Repository[UserInfo], parsers []parser.Parser) *Controller {
	return &Controller{
		client:     client,
		repository: repository,
		parsers:    parsers,
	}
}

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func (c *Controller) Process(event events.Event) error {
	switch event.Type {
	case events.Command:
		{
			return c.processCmd(event)
		}
	case events.Message:
		{
			return c.processMessage(event)
		}
	default:
		{
			return e.Wrap("can't process message", ErrUnknownEventType)
		}
	}
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)

	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func (c *Controller) Fetch(limit int) ([]events.Event, error) {
	updates, err := c.client.FetchUpdates(c.offset, limit)

	if err != nil {
		return nil, e.Wrap("can't get controller", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	c.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func event(u tgclient.Update) events.Event {
	updateType := fetchType(u)

	res := events.Event{
		UserHash:    "",
		Type:        updateType,
		Text:        fetchText(u),
		CommandInfo: fetchCommand(u),
	}

	if updateType == events.Callback {
		log.Printf(u.CallbackQuery.Message.Text)
	}

	if updateType == events.Message || updateType == events.Command {
		res.UserHash = md5Hash(fmt.Sprintf("%s%d", u.Message.From.Username, u.Message.Chat.ID))

		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}

	return res
}

func fetchType(u tgclient.Update) events.EventType {
	if u.Message != nil {
		if u.Message.IsCommand() {
			return events.Command
		}

		return events.Message
	} else if u.CallbackQuery != nil {
		return events.Callback
	}

	return events.Unknown
}

func fetchText(u tgclient.Update) string {
	if u.Message == nil {
		return ""
	}

	return u.Message.Text
}

func fetchCommand(u tgclient.Update) events.CommandInfo {
	if !u.Message.IsCommand() {
		return events.CommandInfo{}
	}

	return events.CommandInfo{
		Command:   u.Message.Command(),
		Arguments: u.Message.CommandArguments(),
	}
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
