package pg

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

type MessageRepo struct {
	db *DB
}

func NewMessageRepo(db *DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (repo *MessageRepo) GetMessageByName(name string) (*model.Message, error) {
	message := &model.Message{}

	rows, err := repo.db.Query("SELECT * FROM message WHERE title=$1;", name)
	if err != nil {
		return nil, &utils.GettingError{Name: "message by name", ErrorValue: err}
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&message.ID, &message.UserID, &message.ChannelID, &message.Title, &message.ImageURL)
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

	row := repo.db.QueryRow(
		"INSERT INTO message(channel_id, user_id, title, message_url, imageurl) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		message.ChannelID, message.UserID, message.Title, message.MessageURL, message.ImageURL,
	)
	if err := row.Scan(&id); err != nil {
		return id, &utils.CreateError{Name: "message", ErrorValue: err}
	}

	return id, nil
}

func (repo *MessageRepo) DeleteMessageByID(messageID int) (int, error) {
	var id int

	row := repo.db.QueryRow("DELETE FROM message WHERE id=$1 RETURNING id;", messageID)
	if err := row.Scan(&id); err != nil {
		return id, fmt.Errorf("error while deleting message: %w", err)
	}

	return id, nil
}

func (repo *MessageRepo) GetMessagesWithRepliesCount() ([]model.Message, error) {
	messages := make([]model.Message, 0)

	rows, err := repo.db.Query("SELECT m.id, COUNT(r.id) FROM message m LEFT JOIN replie r ON r.message_id = m.id GROUP BY m.id ORDER BY m.id;")
	if err != nil {
		return nil, &utils.GettingError{Name: "messages with replies count", ErrorValue: err}
	}

	defer rows.Close()
	for rows.Next() {
		message := model.Message{}
		err := rows.Scan(&message.ID, &message.RepliesCount)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages, nil
}
