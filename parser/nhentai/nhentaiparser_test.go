package nhentai

import (
	"testing"
)

func TestParseAll(t *testing.T) {
	parser := New("nhentai.net")

	mangaList, err := parser.ParseAll("amashiro mio")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedMangaCount := 11
	if len(mangaList) != expectedMangaCount {
		t.Errorf("Expected %d manga, got %d", expectedMangaCount, len(mangaList))
	}
}

func TestParseQuantity(t *testing.T) {
	parser := New("nhentai.net")

	quantity, err := parser.ParseQuantity("amashiro mio")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedQuantity := 11
	if quantity != expectedQuantity {
		t.Errorf("Expected quantity %d, got %d", expectedQuantity, quantity)
	}
}
