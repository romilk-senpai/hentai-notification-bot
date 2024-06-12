package tgclient

import (
	"encoding/json"
	"hentai-notification-bot-re/lib/e"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) FetchUpdates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)

	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't process request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)

	return body, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("parse_mode", "HTML")

	_, err := c.doRequest(sendMessageMethod, q)

	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendStandardMarkup(chatID int) error {
	q := url.Values{}
	q.Add("text", "What to do... What to do...")
	q.Add("chat_id", strconv.Itoa(chatID))

	markup, err := json.Marshal(StandardKeyboardMarkup())

	if err != nil {
		return err
	}

	q.Add("reply_markup", string(markup))

	_, err = c.doRequest(sendMessageMethod, q)

	return err
}

func (c *Client) SendTagManager(chatID int, tagGroups []string) error {
	q := url.Values{}
	q.Add("text", "Manage tags")
	q.Add("chat_id", strconv.Itoa(chatID))

	markup, err := json.Marshal(TagManagerInlineMarkup(tagGroups))

	if err != nil {
		return err
	}

	q.Add("reply_markup", string(markup))

	resp, err := c.doRequest(sendMessageMethod, q)

	log.Printf(string(resp))

	return err
}
