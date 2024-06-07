package jsonrepository

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type MockIdentifiable struct {
	Uuid           string         `json:"uuid"`
	Username       string         `json:"username"`
	ChatID         string         `json:"chat_id"`
	SubscribedTags map[string]int `json:"subscribed_tags"`
}

func (m MockIdentifiable) GetUuid() string {
	return m.Uuid
}

func TestJsonRepository_Create(t *testing.T) {
	dir := t.TempDir()
	repo := JsonRepository[MockIdentifiable]{Path: dir}

	record := MockIdentifiable{
		Uuid:     "uuid",
		Username: "username",
		ChatID:   "chat_id",
		SubscribedTags: map[string]int{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": 912,
		},
	}

	createdRecord, err := repo.Create(record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdRecord.GetUuid() != record.GetUuid() || createdRecord.Username != record.Username ||
		createdRecord.ChatID != record.ChatID || !reflect.DeepEqual(createdRecord.SubscribedTags, record.SubscribedTags) {
		t.Fatalf("expected %v, got %v", record, createdRecord)
	}

	filePath := filepath.Join(dir, "uuid.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fmt.Printf("Content of %s after reading:\n%s\n", filePath, string(data))
}

func TestJsonRepository_Read(t *testing.T) {
	dir := t.TempDir()
	repo := JsonRepository[MockIdentifiable]{Path: dir}

	record := MockIdentifiable{
		Uuid:     "uuid",
		Username: "username",
		ChatID:   "chat_id",
		SubscribedTags: map[string]int{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": 912,
		},
	}

	_, err := repo.Create(record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	readRecord, err := repo.Read(record.GetUuid())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if readRecord.GetUuid() != record.GetUuid() || readRecord.Username != record.Username ||
		readRecord.ChatID != record.ChatID || !reflect.DeepEqual(readRecord.SubscribedTags, record.SubscribedTags) {
		t.Fatalf("expected %v, got %v", record, readRecord)
	}
}

func TestJsonRepository_Update(t *testing.T) {
	dir := t.TempDir()
	repo := JsonRepository[MockIdentifiable]{Path: dir}

	record := MockIdentifiable{
		Uuid:     "uuid",
		Username: "username",
		ChatID:   "chat_id",
		SubscribedTags: map[string]int{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": 912,
		},
	}
	_, err := repo.Create(record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	newRecord := MockIdentifiable{
		Uuid:     "uuid",
		Username: "new_username",
		ChatID:   "new_chat_id", SubscribedTags: map[string]int{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": 912,
		},
	}
	updatedRecord, err := repo.Update(record.GetUuid(), newRecord)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedRecord.GetUuid() != newRecord.GetUuid() || updatedRecord.Username != newRecord.Username ||
		updatedRecord.ChatID != newRecord.ChatID || !reflect.DeepEqual(updatedRecord.SubscribedTags, newRecord.SubscribedTags) {
		t.Fatalf("expected %v, got %v", newRecord, updatedRecord)
	}
}

func TestJsonRepository_Delete(t *testing.T) {
	dir := t.TempDir()
	repo := JsonRepository[MockIdentifiable]{Path: dir}

	record := MockIdentifiable{
		Uuid:     "uuid",
		Username: "username",
		ChatID:   "chat_id",
		SubscribedTags: map[string]int{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": 912,
		},
	}
	_, err := repo.Create(record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = repo.Delete(record.GetUuid())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.Exists(record.GetUuid()) {
		t.Fatalf("expected record to be deleted, but it still exists")
	}
}

func TestJsonRepository_Exists(t *testing.T) {
	dir := t.TempDir()
	repo := JsonRepository[MockIdentifiable]{Path: dir}

	record := MockIdentifiable{
		Uuid:     "uuid",
		Username: "username",
		ChatID:   "chat_id",
		SubscribedTags: map[string]int{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": 912,
		},
	}
	_, err := repo.Create(record)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repo.Exists(record.GetUuid()) {
		t.Fatalf("expected record to exist, but it doesn't")
	}

	if repo.Exists("non_existing_uuid") {
		t.Fatalf("expected record not to exist, but it does")
	}
}
