package pg

import (
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

type ChannelPgRepo struct {
	db *DB
}

func NewChannelRepo(db *DB) *ChannelPgRepo {
	return &ChannelPgRepo{db}
}

func (repo *ChannelPgRepo) GetChannelByName(name string) (*model.Channel, error) {
	channel := &model.Channel{}

	rows, err := repo.db.Query("SELECT * FROM channel WHERE name=$1;", name)
	if err != nil {
		return nil, &utils.GettingError{Name: "channel by name", ErrorValue: err}
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
		return id, &utils.CreateError{Name: "channel", ErrorValue: err}
	}

	return id, nil
}
