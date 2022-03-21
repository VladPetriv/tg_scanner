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

func (repo *ChannelPgRepo) GetChannels() (*[]model.Channel, error) {
	chanels := make([]model.Channel, 0)
	rows, err := repo.db.Query("SELECT * FROM channels;")
	if err != nil {
		return nil, fmt.Errorf("Error while getting channels: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		channel := model.Channel{}
		err := rows.Scan(&channel.Id, &channel.Name)
		if err != nil {
			continue
		}

		chanels = append(chanels, channel)
	}

	return &chanels, nil
}

func (repo *ChannelPgRepo) GetChannel(channelId int) (*model.Channel, error) {
	channel := model.Channel{}

	rows, err := repo.db.Query("SELECT * FROM channels;")
	if err != nil {
		return nil, fmt.Errorf("Error while getting channel: %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&channel.Id, &channel.Name)
		if err != nil {
			continue
		}
	}

	return &channel, nil
}

func (repo *ChannelPgRepo) CreateChannel(channel *model.Channel) (*model.Channel, error) {
	_, err := repo.db.Exec("INSERT INTO channel(name) VALUES ($1)", channel.Name)
	if err != nil {
		return nil, fmt.Errorf("Error while creating channel: %w", err)
	}

	return channel, nil
}

func (repo *ChannelPgRepo) DeleteChannel(channelId int) error {
	_, err := repo.db.Exec("DELETE FROM channels WHERE id = $1;", channelId)
	if err != nil {
		return fmt.Errorf("Error while deleting channel: %w", err)
	}

	return nil
}
