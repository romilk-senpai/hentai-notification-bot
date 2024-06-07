package nhentai

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"hentai-notification-bot-re/lib/e"
	"hentai-notification-bot-re/parser"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	cookie    = "cf_clearance=URomM.uxgeqeCA5WIeTkjZvZDOWgCyqS0x9pGl5X4_0-1717624444-1.0.1.1-GGJX6eqEJwVLCEokc1c6k7RraeK94u5_OUjkm8.xMaj1MXlUShSJgt13dYlOF2VXbFX5.VlbauUE3QnRYB7Chw; csrftoken=A8Evs3ba7wcJWYAwMSr8yV9b146eHrnF4aWhPGVrJ92kHOKP1xuQIpmL1qYr5Oii;"
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

func (p *Parser) ParseOne(query string) (*parser.Manga, error) {
	return nil, nil
}

func (p *Parser) ParseAll(query string) (manga []parser.Manga, err error) {
	defer func() { err = e.WrapIfErr("can't process request", err) }()

	data, err := p.doRequest("/search/", "q="+strings.ReplaceAll(query, " ", "+"))

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

func (p *Parser) doRequest(path string, rawQuery string) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't process request", err) }()

	u := url.URL{
		Scheme:   "https",
		Host:     p.host,
		Path:     path,
		RawQuery: rawQuery,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Cookie", cookie)

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
