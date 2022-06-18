package pg

import (
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
	replie := &model.Replie{}

	rows, err := repo.db.Query("SELECT * FROM replie WHERE title=$1;", name)
	if err != nil {
		return nil, &utils.GettingError{Name: "replie by name", ErrorValue: err}
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&replie.ID, &replie.UserID, &replie.MessageID, &replie.Title)
		if err != nil {
			continue
		}
	}

	if replie.Title == "" {
		return nil, nil
	}

	return replie, nil
}

func (repo *ReplieRepo) CreateReplie(replie *model.Replie) (int, error) {
	var id int

	row := repo.db.QueryRow(
		"INSERT INTO replie (user_id, message_id, title) VALUES ($1, $2, $3) RETURNING id;",
		replie.UserID, replie.MessageID, replie.Title,
	)
	if err := row.Scan(&id); err != nil {
		return id, &utils.CreateError{Name: "replie", ErrorValue: err}
	}

	return id, nil
}
