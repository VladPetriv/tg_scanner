package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/errors"
	"github.com/VladPetriv/tg_scanner/pkg/filter"
)

func WriteMessagesToFile(messages []model.TgMessage, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("error while open file: %w", err)
	}

	encodedMessages, err := json.Marshal(messages)
	if err != nil {
		return &errors.CreateError{Name: "JSON", ErrorValue: err}
	}

	_, err = file.Write(encodedMessages)
	if err != nil {
		return fmt.Errorf("error while writing to file: %w", err)
	}

	return nil
}

func GetMessagesFromFile(fileName string) ([]model.TgMessage, error) {
	messages := make([]model.TgMessage, 0, 10)

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("read file error: %w", err)
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, &errors.CreateError{Name: "JSON", ErrorValue: err}
	}

	file, err := os.Create(fileName)
	if err != nil {
		return nil, &errors.CreateError{Name: "file", ErrorValue: err}
	}

	_, err = file.WriteString("")
	if err != nil {
		return nil, fmt.Errorf("write file error: %w", err)
	}

	return messages, nil
}

func CreateFilesForGroups(groups []model.TgGroup) error {
	for _, group := range groups {
		fileName := fmt.Sprintf("%s.json", group.Username)
		if _, err := os.Stat("./data/" + fileName); err == nil {
			continue
		}

		file, err := os.Create(fileName)
		if err != nil {
			return &errors.CreateError{Name: "file", ErrorValue: err}
		}

		_, err = file.WriteString("[ ]")
		if err != nil {
			return fmt.Errorf("write file error: %s", err)
		}

		err = os.Rename(fileName, fmt.Sprintf("./data/%s", fileName))
		if err != nil {
			return fmt.Errorf("rename file error:%s", err)
		}
	}

	return nil
}

func CreateFileForIncoming() error {
	file, err := os.Create("./data/incoming.json")
	if os.IsNotExist(err) {
		return &errors.CreateError{Name: "file", ErrorValue: err}
	}

	_, err = file.WriteString("[ ]")
	if err != nil {
		return fmt.Errorf("write file error: %w", err)
	}

	return nil
}

func InitDirectories() error {
	dirs := [3]string{"data", "images"}

	for _, dir := range dirs {
		err := os.Mkdir(dir, 0o755)
		if os.IsNotExist(err) {
			return &errors.CreateError{Name: "directory", ErrorValue: err}
		}
	}

	return nil
}

func ParseFromFiles(path string) ([]model.TgMessage, error) {
	messages := make([]model.TgMessage, 0, 10)

	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open dir error: %w", err)
	}

	files, err := dir.ReadDir(0)
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}

	for _, file := range files {
		pathToFile := fmt.Sprintf("./%s/%s", path, file.Name())

		data, err := GetMessagesFromFile(pathToFile)
		if err != nil {
			fmt.Printf("%s - %s\n", file.Name(), err)

			continue
		}

		err = os.Remove(pathToFile)
		if err != nil {
			return nil, fmt.Errorf("remove file error: %w", err)
		}

		f, err := os.Create(pathToFile)
		if err != nil {
			return nil, fmt.Errorf("open file error: %w", err)
		}

		_, err = f.WriteString("[  ]")
		if err != nil {
			return nil, fmt.Errorf("write file error: %w", err)
		}

		messages = append(messages, data...)
	}

	result := filter.RemoveDuplicatesFromMessages(messages)

	return result, nil
}
