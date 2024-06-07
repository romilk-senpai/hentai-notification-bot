package tgcontroller

import (
	"errors"
	"fmt"
	tgclient "hentai-notification-bot-re/client/telegram"
	"hentai-notification-bot-re/controller"
	"hentai-notification-bot-re/lib/e"
	"hentai-notification-bot-re/parser"
	"hentai-notification-bot-re/repository"
	"log"
	"strings"
)

type Controller struct {
	client     *tgclient.Client
	offset     int
	repository *repository.Repository[any]
	parsers    []parser.Parser
}

func New(client *tgclient.Client, repository *repository.Repository[any], parsers []parser.Parser) *Controller {
	return &Controller{
		client:     client,
		repository: repository,
		parsers:    parsers,
	}
}

type Meta struct {
	ChatID   int
	Username string
}

var ErrUnknownEventType = errors.New("unknown event type")
var ErrUnknownMetaType = errors.New("unknown meta type")

func (p *Controller) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		{
			return p.processMessage(event)
		}
	default:
		{
			return e.Wrap("can't process message", ErrUnknownEventType)
		}
	}
}

func (p *Controller) processMessage(event events.Event) error {
	meta, err := meta(event)

	if err != nil {
		return e.Wrap("can't process message", err)
	}

	log.Printf("message from %s: %s", meta.Username, event.Text)

	if len(event.Text) == 0 {
		return nil
	}

	for _, mParser := range p.parsers {
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

		if err := p.client.SendMessage(meta.ChatID, responseBuilder.String()); err != nil {
			return e.Wrap("can't send message", err)
		}
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)

	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func (p *Controller) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.client.FetchUpdates(p.offset, limit)

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

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func event(u tgclient.Update) events.Event {
	updateType := fetchType(u)

	res := events.Event{
		Type: updateType,
		Text: fetchText(u),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}

	return res
}

func fetchType(u tgclient.Update) events.EventType {
	if u.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(u tgclient.Update) string {
	if u.Message == nil {
		return ""
	}

	return u.Message.Text
}
