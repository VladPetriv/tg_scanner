package pg

import (
	"database/sql"

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
	channel := model.Channel{}

	err := repo.db.Get(&channel, "SELECT * FROM channel WHERE name=$1;", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, &utils.GettingError{Name: "channel by name", ErrorValue: err}
	}

	return &channel, nil
}

func (repo *ChannelPgRepo) CreateChannel(channel *model.Channel) (int, error) {
	var id int

	row := repo.db.QueryRow("INSERT INTO channel(name, title, imageurl) VALUES ($1, $2, $3) RETURNING id;", channel.Name, channel.Title, channel.ImageURL)

	if err := row.Scan(&id); err != nil {
		return id, &utils.CreateError{Name: "channel", ErrorValue: err}
	}

	return id, nil
}
