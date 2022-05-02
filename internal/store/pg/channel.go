package pg

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type ChannelPgRepo struct {
	db *DB
}

func NewChannelRepo(db *DB) *ChannelPgRepo {
	return &ChannelPgRepo{db}
}

func (repo *ChannelPgRepo) GetChannel(channelID int) (*model.Channel, error) {
	channel := &model.Channel{}

	rows, err := repo.db.Query("SELECT * FROM channel WHERE id=$1;", channelID)
	if err != nil {
		return nil, fmt.Errorf("error while getting channel: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&channel.ID, &channel.Name, &channel.Title, &channel.PhotoURL)
		if err != nil {
			continue
		}
	}

	if channel.Name == "" {
		return nil, fmt.Errorf("channel not found")
	}

	return channel, nil
}

func (repo *ChannelPgRepo) GetChannelByName(name string) (*model.Channel, error) {
	channel := &model.Channel{}

	rows, err := repo.db.Query("SELECT * FROM channel WHERE name=$1;", name)
	if err != nil {
		return nil, fmt.Errorf("error while getting channel: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&channel.ID, &channel.Name, &channel.Title, &channel.PhotoURL)
		if err != nil {
			continue
		}
	}

	if channel.Name == "" {
		return nil, nil
	}

	return channel, nil
}

func (repo *ChannelPgRepo) CreateChannel(channel *model.Channel) (int, error) {
	var id int
	row := repo.db.QueryRow("INSERT INTO channel(name, title, photourl) VALUES ($1, $2, $3) RETURNING id;", channel.Name, channel.Title, channel.PhotoURL)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while creating channel: %w", err)
	}

	return 1, nil
}
