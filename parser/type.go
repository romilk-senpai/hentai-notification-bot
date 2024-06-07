package parser

type Parser interface {
	ParseOne(query string) (*Manga, error)
	ParseAll(query string) ([]Manga, error)
}

type Manga struct {
	Name string
	Url  string
}
