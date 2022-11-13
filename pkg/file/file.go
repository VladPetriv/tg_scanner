package file

import (
	"fmt"
	"os"
)

func CreateFileForIncoming() error {
	file, err := os.Create("./data/incoming.json")
	if os.IsNotExist(err) {
		return fmt.Errorf("create file error: %w", err)
	}

	_, err = file.WriteString("[ ]")
	if err != nil {
		return fmt.Errorf("write to file error: %w", err)
	}

	return nil
}

func InitDirectories() error {
	dirs := [2]string{"data", "images"}

	for _, dir := range dirs {
		err := os.Mkdir(dir, 0o755)
		if os.IsNotExist(err) {
			return fmt.Errorf("create directory error: %w", err)
		}
	}

	return nil
}
