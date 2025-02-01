package nhentai

import (
	"bytes"
	"errors"
	"hentai-notification-bot-re/lib/e"
	"hentai-notification-bot-re/lib/e/config"
	"hentai-notification-bot-re/parser"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

const (
	ParserName = "nhentai"
)

type Parser struct {
	host   string
	client http.Client
	cfg    *config.Config
}

func New(host string, cfg *config.Config) *Parser {
	return &Parser{
		host:   host,
		client: http.Client{},
		cfg:    cfg,
	}
}

func (p *Parser) ParserName() string {
	return ParserName
}

func (p *Parser) ParseOne(query string) (*parser.Manga, error) {
	return nil, nil
}

func (p *Parser) ParseAll(query string) (manga []parser.Manga, err error) {
	defer func() { err = e.WrapIfErr("can't process request", err) }()

	data, err := p.doRequest("/search/", "q="+url.QueryEscape(query))

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return nil, err
	}

	mangoes := make([]parser.Manga, 0)

	contentBlock := doc.Find("div#content")
	containerBlock := contentBlock.Find("div.container.index-container")

	containerBlock.Find("div.gallery").Each(func(i int, selection *goquery.Selection) {
		mangaHref, _ := selection.Find("a.cover").First().Attr("href")
		mangaName := selection.Find("div.caption").First().Text()

		mangaUrl := url.URL{
			Scheme: "https",
			Host:   p.host,
			Path:   mangaHref,
		}

		manga := parser.Manga{
			Name: mangaName,
			Url:  mangaUrl.String(),
		}

		mangoes = append(mangoes, manga)
	})

	return mangoes, nil
}

func (p *Parser) ParseQuantity(query string) (quantity int, err error) {
	defer func() { err = e.WrapIfErr("can't process request", err) }()

	data, err := p.doRequest("/search/", "q="+strings.ReplaceAll(query, " ", "+"))

	if err != nil {
		return 0, err
	}

	reader := bytes.NewReader(data)

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return 0, err
	}

	contentBlock := doc.Find("div#content")
	resultCountEl := contentBlock.Find("h1").First().Text()

	return parseNumeric(resultCountEl)
}

func (p *Parser) QueryToLink(query string) string {
	u := url.URL{
		Scheme:   "https",
		Host:     p.host,
		Path:     "/search/",
		RawQuery: "q=" + url.QueryEscape(query),
	}

	return u.String()
}

func parseNumeric(input string) (int, error) {
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, input)

	output, err := strconv.Atoi(cleaned)

	if err != nil {
		return 0, err
	}

	return output, nil
}

func (p *Parser) doRequest(path string, rawQuery string) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't process request", err) }()

	u := url.URL{
		Scheme:   "https",
		Host:     p.host,
		Path:     path,
		RawQuery: rawQuery,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	req.Header.Add("User-Agent", p.cfg.UserAgent)
	req.Header.Add("Cookie", p.cfg.Nhcookie)

	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		log.Printf("%s returned %d", u.String(), resp.StatusCode)
		return nil, errors.New("parser request error")
	}

	body, err := io.ReadAll(resp.Body)

	return body, err
}
