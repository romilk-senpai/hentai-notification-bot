package hentaifox

import (
	"testing"
)

func TestParseAll(t *testing.T) {
	parser := New("hentaifox.com")

	mangaList, err := parser.ParseAll("itohana")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedMangaCount := 7
	if len(mangaList) != expectedMangaCount {
		t.Errorf("Expected %d manga, got %d", expectedMangaCount, len(mangaList))
	}
}

func TestParseQuantity(t *testing.T) {
	parser := New("hentaifox.com")

	quantity, err := parser.ParseQuantity("itohana")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedQuantity := 7
	if quantity != expectedQuantity {
		t.Errorf("Expected quantity %d, got %d", expectedQuantity, quantity)
	}
}
