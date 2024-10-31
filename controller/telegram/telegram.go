package tgcontroller

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	tgclient "hentai-notification-bot-re/client/telegram"
	events "hentai-notification-bot-re/controller"
	"hentai-notification-bot-re/lib/e"
	"hentai-notification-bot-re/parser"
	"hentai-notification-bot-re/repository"
	"log"
	"net/http"
)

type Controller struct {
	client     *tgclient.Client
	offset     int
	repository repository.Repository[*UserInfo]
	parsers    []parser.Parser
}

func New(client *tgclient.Client, repository repository.Repository[*UserInfo], parsers []parser.Parser) *Controller {
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
	case events.Callback:
		{
			return c.processCallback(event)
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

func NewHandler(p events.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()

		var res tgclient.Update

		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			log.Printf("bad requset: %s", err.Error())
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		updates := []tgclient.Update{res}

		events, err := getEvents(updates)

		if err != nil {
			log.Printf("can't handle event: %s", err.Error())
			return
		}

		for _, event := range events {
			log.Printf("got new event: %v", event.Type)

			if err := p.Process(event); err != nil {
				log.Printf("can't handle event: %s", err.Error())

				continue
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request processed successfully"))
	}
}

func (c *Controller) Fetch(limit int) ([]events.Event, error) {
	updates, err := c.client.FetchUpdates(c.offset, limit)

	if err != nil {
		return nil, e.Wrap("can't get controller", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	c.offset = updates[len(updates)-1].ID + 1

	return getEvents(updates)
}

func getEvents(updates []tgclient.Update) ([]events.Event, error) {

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

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

	if updateType == events.Message || updateType == events.Command {
		res.UserHash = md5Hash(fmt.Sprintf("%d%d", u.Message.From.ID, u.Message.Chat.ID))

		res.Meta = Meta{
			Update: u,
		}
	} else if updateType == events.Callback {
		res.UserHash = md5Hash(fmt.Sprintf("%d%d", u.CallbackQuery.From.ID, u.CallbackQuery.Message.Chat.ID))

		res.Meta = Meta{
			Update: u,
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
