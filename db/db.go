package db

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/rikughi/technical-test/helper"
	"github.com/rikughi/technical-test/model"
)

var Collection model.Collection

var (
	ErrCodeNotFound = errors.New("code does not exist")
)

func Init() {
	file, err := os.Open("data.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		helper.PanicIfError(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&Collection)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
		helper.PanicIfError(err)
	}
}

func FirstOrNil(code string) *model.Data {
	for _, item := range Collection.Data {
		if item.Code == code {
			return item
		}
	}

	return nil
}

func UpdateByCode(code string, data *model.Data) error {
	for _, item := range Collection.Data {
		if item.Code == code {
			item = data
			return nil
		}
	}

	return ErrCodeNotFound
}

func DeleteByCode(code string) error {
	for i, item := range Collection.Data {
		if item.Code == code {
			Collection.Data = append(Collection.Data[:i], Collection.Data[i+1:]...)
			return nil
		}
	}

	return ErrCodeNotFound
}

func Contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
