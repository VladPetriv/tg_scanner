package pg

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type MessageRepo struct {
	db *DB
}

func NewMessageRepo(db *DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (repo *MessageRepo) GetMessages() ([]model.Message, error) {
	messages := make([]model.Message, 0)

	rows, err := repo.db.Query("SELECT * FROM message;")
	if err != nil {
		return nil, fmt.Errorf("error while getting messages: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		message := model.Message{}
		err := rows.Scan(&message.ID, &message.UserID, &message.ChannelID, &message.Title)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("messages not found")
	}

	return messages, nil
}

func (repo *MessageRepo) GetMessage(messageID int) (*model.Message, error) {
	message := &model.Message{}

	rows, err := repo.db.Query("SELECT * FROM message WHERE id=$1;", messageID)
	if err != nil {
		return nil, fmt.Errorf("error while getting message: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&message.ID, &message.UserID, &message.ChannelID, &message.Title)
		if err != nil {
			continue
		}
	}

	if message.Title == "" {
		return nil, fmt.Errorf("message not found")
	}

	return message, nil
}

func (repo *MessageRepo) GetMessageByName(name string) (*model.Message, error) {
	message := &model.Message{}

	rows, err := repo.db.Query("SELECT * FROM message WHERE title=$1;", name)
	if err != nil {
		return nil, fmt.Errorf("error while getting message: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&message.ID, &message.UserID, &message.ChannelID, &message.Title)
		if err != nil {
			continue
		}
	}

	if message.Title == "" {
		return nil, nil
	}

	return message, nil
}

func (repo *MessageRepo) CreateMessage(message *model.Message) (int, error) {
	var id int
	row := repo.db.QueryRow("INSERT INTO message(channel_id, user_id, title) VALUES ($1, $2, $3) RETURNING id;", message.ChannelID, message.UserID, message.Title)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while creating message: %w", err)
	}

	return id, nil
}

func (repo *MessageRepo) DeleteMessage(messageID int) error {
	_, err := repo.db.Exec("DELETE FROM message WHERE id=$1;", messageID)
	if err != nil {
		return fmt.Errorf("error while deleting message: %w", err)
	}

	return nil
}
