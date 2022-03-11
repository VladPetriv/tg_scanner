package file

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/VladPetriv/tg_scanner/internal/channel"
	"github.com/VladPetriv/tg_scanner/internal/message"
	"github.com/sirupsen/logrus"
)

func WriteMessagesToFile(msgs []message.Message, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644) // nolint
	if err != nil {
		return fmt.Errorf("ERROR_WHILE_OPENING_FILE: %w", err)
	}
	defer file.Close()

	messages, err := json.Marshal(msgs)
	if err != nil {
		return fmt.Errorf("ERROR_WHILE_CREATEING_JSON: %w", err)
	}

	_, err = file.WriteString(string(messages))
	if err != nil {
		return fmt.Errorf("ERROR_WHILE_WRITING_TO_FILE:%w", err)
	}

	return nil
}

func GetMessagesFromFile(fileName string) ([]message.Message, error) {
	var messages []message.Message

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_OPENING_FILE: %w", err)
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_CREATEING_JSON: %w", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_CREATING_FILE:%w", err)
	}

	_, err = file.WriteString("")
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_WRITING_TO_FILE:%w", err)
	}

	return messages, nil
}

func CreateFilesForGroups(groups []channel.Group) {
	var once sync.Once

	once.Do(func() {
		for _, group := range groups {
			fileName := fmt.Sprintf("%s.json", group.Username)
			file, err := os.Create(fileName)
			if err != nil {
				logrus.Errorf("ERROR_WHILE_WORKING_WITH_FILES:%s", err)
			}
			_, err = file.WriteString("[]")
			if err != nil {
				logrus.Errorf("ERROR_WHILE_WRITING_TO_FILE:%s", err)
			}

			err = os.Rename(fileName, fmt.Sprintf("./data/%s", fileName))
			if err != nil {
				logrus.Errorf("ERROR_WHILE_WORKING_WITH_FILES:%s", err)
			}
		}
		logrus.Info("Files was created")
	})
}

func CreateFileForIncoming() error {
	file, err := os.Create("./data/incoming.json")
	if os.IsNotExist(err) {
		return fmt.Errorf("ERROR_WHILE_CREATE_FILE:%w", err)
	}

	_, err = file.WriteString("[]")
	if err != nil {
		return fmt.Errorf("ERROR_WHILE_WRITING_TO_FILE:%w", err)
	}

	return nil
}

func CreateFileForLogger(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640) // nolint
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_OPEN_FILE:%w", err)
	}

	return file, nil
}

func CreateDirs() error {
	err := os.Mkdir("data", 0o755) // nolint
	if os.IsNotExist(err) {
		return fmt.Errorf("ERROR_WHILE_CREATE_DIR:%w", err)
	}

	err = os.Mkdir("logs", 0o755) // nolint
	if os.IsNotExist(err) {
		return fmt.Errorf("ERROR_WHILE_CREATE_DIR:%w", err)
	}

	return nil
}
