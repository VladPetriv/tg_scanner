package file

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/VladPetriv/tg_scanner/internal/channel"
	"github.com/VladPetriv/tg_scanner/internal/message"
)

var once sync.Once

func WriteMessagesToFile(msgs []message.Message, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ERROR_WHILE_OPENING_FILE: %s", err)
	}
	defer file.Close()

	messages, err := json.Marshal(msgs)
	if err != nil {
		return fmt.Errorf("ERROR_WHILE_CREATEING_JSON: %s", err)
	}

	file.WriteString(string(messages))

	return nil
}

func GetMessagesFromFile(fileName string) ([]message.Message, error) {
	var messages []message.Message
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_OPENING_FILE: %s", err)
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_CREATEING_JSON: %s", err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	file.WriteString("")

	return messages, nil
}

func GetGroupsFromFile(fileName string) ([]channel.Group, error) {
	var gropus []channel.Group
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_OPENING_FILE: %s", err)
	}
	err = json.Unmarshal(data, &gropus)
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_CREATEING_JSON: %s", err)
	}

	return gropus, nil
}

func CreateFiles(groups []channel.Group) {
	once.Do(func() {
		os.Mkdir("data", 0644)
		for _, group := range groups {
			fileName := fmt.Sprintf("%s.json", group.Username)
			file, err := os.Create(fileName)
			if err != nil {
				log.Fatal(err)
			}
			file.WriteString("[]")
		}
		log.Println("File was created")
	})

}
