package nhentai

import (
	"hentai-notification-bot-re/lib/e/config"
	"log"
	"testing"
)

func TestParseAll(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config error")
	}
	parser := New("nhentai.net", cfg)

	mangaList, err := parser.ParseAll("amashiro mio")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedMangaCount := 16
	if len(mangaList) != expectedMangaCount {
		t.Errorf("Expected %d manga, got %d", expectedMangaCount, len(mangaList))
	}
}

func TestParseQuantity(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config error")
	}
	parser := New("nhentai.net", cfg)

	quantity, err := parser.ParseQuantity("amashiro mio")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedQuantity := 16
	if quantity != expectedQuantity {
		t.Errorf("Expected quantity %d, got %d", expectedQuantity, quantity)
	}
}
