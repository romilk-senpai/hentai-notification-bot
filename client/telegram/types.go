package tgclient

import "strings"

const (
	GetMangaUpdatesMessage = "Get manga updates"
	ManageTagsMessage      = "Manage tags"
)

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message,omitempty"`
	CallbackQuery *CallbackQuery   `json:"callback_query,omitempty"`
}

type CallbackQuery struct {
	ID      string           `json:"id"`
	From    *User            `json:"from"`
	Message *IncomingMessage `json:"message,omitempty"`
	Data    string           `json:"data,omitempty"`
}

func (callbackData *CallbackQuery) ParseCallbackData() *CallbackData {
	before, after, found := strings.Cut(callbackData.Data, ":")
	data := CallbackData{}
	data.Key = before
	if found {
		data.Value = after
	} else {
		data.Value = before
	}
	return &data
}

type CallbackData struct {
	Key   string
	Value string
}

type SendMessageResponse struct {
	Ok     bool                      `json:"ok"`
	Result SendMessageResponseResult `json:"result"`
}

type SendMessageResponseResult struct {
	MessageID int `json:"message_id"`
}

type IncomingMessage struct {
	Text     string          `json:"text"`
	ID       int             `json:"message_id"`
	From     User            `json:"from"`
	Chat     Chat            `json:"chat"`
	Entities []MessageEntity `json:"entities"`
}

func (m *IncomingMessage) IsCommand() bool {
	if m == nil || m.Entities == nil || len(m.Entities) == 0 {
		return false
	}

	entity := m.Entities[0]
	return entity.Offset == 0 && entity.IsCommand()
}

func (m *IncomingMessage) Command() string {
	if !m.IsCommand() {
		return ""
	}

	entity := m.Entities[0]
	return m.Text[1:entity.Length]
}

func (m *IncomingMessage) CommandArguments() string {
	if !m.IsCommand() {
		return ""
	}

	entity := m.Entities[0]

	if len(m.Text) == entity.Length {
		return "" // The command makes up the whole message
	}

	return m.Text[entity.Length+1:]
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

func (e *MessageEntity) IsCommand() bool {
	return e.Type == "bot_command"
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type InlineKeyboard struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

func TagManagerInlineMarkup(tagGroups []string) InlineKeyboard {
	keyboard := [][]InlineKeyboardButton{
		{
			{Text: "Add tag group", CallbackData: "addTagGroup"},
		},
	}

	for _, tagGroup := range tagGroups {
		row := []InlineKeyboardButton{
			{Text: tagGroup, CallbackData: "nothing"},
			{Text: "❌", CallbackData: "deleteTags:" + tagGroup},
		}

		keyboard = append(keyboard, row)
	}

	cancelRow := []InlineKeyboardButton{
		{Text: "Cancel", CallbackData: "cancelManage"},
	}

	keyboard = append(keyboard, cancelRow)

	return InlineKeyboard{
		InlineKeyboard: keyboard,
	}
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

func StandardKeyboardMarkup() ReplyKeyboardMarkup {
	return ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{{Text: GetMangaUpdatesMessage}, {Text: ManageTagsMessage}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}
