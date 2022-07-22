package pg

import (
	"database/sql"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

type ReplieRepo struct {
	db *DB
}

func NewReplieRepo(db *DB) *ReplieRepo {
	return &ReplieRepo{db: db}
}

func (repo *ReplieRepo) GetReplieByName(name string) (*model.Replie, error) {
	replie := model.Replie{}

	err := repo.db.Get(&replie, "SELECT * FROM replie WHERE title=$1;", name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, &utils.GettingError{Name: "replie by name", ErrorValue: err}
	}

	return &replie, nil
}

func (repo *ReplieRepo) CreateReplie(replie *model.Replie) (int, error) {
	var id int

	row := repo.db.QueryRow(
		"INSERT INTO replie (user_id, message_id, title, imageurl) VALUES ($1, $2, $3, $4) RETURNING id;",
		replie.UserID, replie.MessageID, replie.Title, replie.ImageURL,
	)
	if err := row.Scan(&id); err != nil {
		return id, &utils.CreateError{Name: "replie", ErrorValue: err}
	}

	return id, nil
}
