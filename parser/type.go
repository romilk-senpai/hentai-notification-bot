package parser

type Parser interface {
	ParserName() string
	ParseOne(query string) (*Manga, error)
	ParseAll(query string) ([]Manga, error)
	ParseQuantity(query string) (int, error)
}

type Manga struct {
	Name string
	Url  string
}
