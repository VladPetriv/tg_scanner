package pg

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type MessageRepo struct {
	db *DB
}

func NewMessageRepo(db *DB) *MessageRepo {
	return &MessageRepo{db}
}

func (repo *MessageRepo) GetMessages() (*[]model.Message, error) {
	messages := make([]model.Message, 0)

	rows, err := repo.db.Query("SELECT * FROM message;")
	if err != nil {
		return nil, fmt.Errorf("error while getting messages: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		message := model.Message{}
		err := rows.Scan(&message.ID, &message.ChannelID, &message.Title)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	return &messages, nil
}

func (repo *MessageRepo) GetMessage(messageId int) (*model.Message, error) {
	message := &model.Message{}

	rows, err := repo.db.Query("SELECT * FROM message WHERE id=$1", messageId)
	if err != nil {
		return nil, fmt.Errorf("error while getting message: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&message.ID, &message.ChannelID, &message.Title)
		if err != nil {
			continue
		}
	}

	return message, nil
}

func (repo *MessageRepo) GetMessageByName(name string) (*model.Message, error) {
	message := &model.Message{}

	rows, err := repo.db.Query("SELECT * FROM message WHERE title=$1", name)
	if err != nil {
		return nil, fmt.Errorf("error while getting message: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&message.ID, &message.ChannelID, &message.Title)
		if err != nil {
			continue
		}
	}

	return message, nil
}

func (repo *MessageRepo) CreateMessage(message *model.Message) error {
	_, err := repo.db.Exec("INSERT INTO message(channel_id,title) VALUES ($1,$2)", message.ChannelID, message.Title)
	if err != nil {
		return fmt.Errorf("error while creating message: %w", err)
	}

	return nil
}

func (repo *MessageRepo) DeleteMessage(messageId int) error {
	_, err := repo.db.Exec("DELETE FROM message WHERE id=$1", messageId)
	if err != nil {
		return fmt.Errorf("error while deleting message: %w", err)
	}

	return nil
}
