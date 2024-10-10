package hentaifox

import (
	"bytes"
	"errors"
	"hentai-notification-bot-re/lib/e"
	"hentai-notification-bot-re/parser"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	ParserName = "hentaifox"
)

type Parser struct {
	host   string
	client http.Client
}

func New(host string) *Parser {
	return &Parser{
		host:   host,
		client: http.Client{},
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

	galleryBlock := doc.Find("div.lc_galleries")

	galleryBlock.Find("div.thumb").Each(func(i int, selection *goquery.Selection) {
		mangaHref, _ := selection.Find("div.inner_thumb").Find("a").First().Attr("href")
		mangaName := selection.Find("div.caption").Find("a").First().Text()

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

	overviewBlock := doc.Find("div.galleries_overview.g_center")
	resultCountEl := overviewBlock.Find("h1").First().Text()

	return parser.ParseNumeric(resultCountEl)
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

func (p *Parser) doRequest(path string, rawQuery string) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't process request", err) }()

	u := url.URL{
		Scheme:   "https",
		Host:     p.host,
		Path:     path,
		RawQuery: rawQuery,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

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
