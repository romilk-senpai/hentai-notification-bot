package jsonrepository

import (
	"encoding/json"
	"fmt"
	"hentai-notification-bot-re/repository"
	"log"
	"os"
	"path/filepath"
)

type JsonRepository[T repository.Identifiable] struct {
	Path string
}

const defaultPerm = 0774

func (r *JsonRepository[T]) New(path string) *JsonRepository[T] {
	return &JsonRepository[T]{
		Path: path,
	}
}

func (r *JsonRepository[T]) Create(record T) (T, error) {
	if err := os.MkdirAll(r.Path, defaultPerm); err != nil {
		return record, err
	}

	filePath := filepath.Join(r.Path, record.GetUuid()+".json")

	file, err := os.Create(filePath)

	if err != nil {
		return record, err
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Print(err.Error())
		}
	}(file)

	data, err := json.Marshal(record)

	if err != nil {
		return record, err
	}

	_, err = file.Write(data)

	if err != nil {
		return record, err
	}

	return record, nil
}

func (r *JsonRepository[T]) Read(uuid string) (T, error) {
	var record T

	filePath := filepath.Join(r.Path, uuid+".json")

	data, err := os.ReadFile(filePath)

	if err != nil {
		return record, err
	}

	err = json.Unmarshal(data, &record)

	if err != nil {
		return record, err
	}

	return record, nil
}

func (r *JsonRepository[T]) Update(uuid string, newRecord T) (T, error) {
	exists := r.Exists(uuid)

	if !exists {
		return newRecord, fmt.Errorf("%s record does not exits", uuid)
	}

	_, err := r.Create(newRecord)

	if err != nil {
		return newRecord, err
	}

	return newRecord, nil
}

func (r *JsonRepository[T]) Delete(uuid string) error {
	exists := r.Exists(uuid)

	if !exists {
		return fmt.Errorf("%s record does not exits", uuid)
	}

	filePath := filepath.Join(r.Path, uuid+".json")

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func (r *JsonRepository[T]) Exists(uuid string) bool {
	filePath := filepath.Join(r.Path, uuid+".json")

	_, err := os.Stat(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		return false
	}

	return true
}
