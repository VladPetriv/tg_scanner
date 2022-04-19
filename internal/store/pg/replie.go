package pg

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type ReplieRepo struct {
	db *DB
}

func NewReplieRepo(db *DB) *ReplieRepo {
	return &ReplieRepo{db: db}
}

func (repo *ReplieRepo) GetReplies() ([]model.Replie, error) {
	replies := make([]model.Replie, 0)
	rows, err := repo.db.Query("SELECT * FROM replie;")
	if err != nil {
		return nil, fmt.Errorf("error while getting replies: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		replie := model.Replie{}
		err := rows.Scan(&replie.ID, &replie.UserID, &replie.MessageID, &replie.Title)
		if err != nil {
			continue
		}

		replies = append(replies, replie)
	}

	if len(replies) == 0 {
		return nil, fmt.Errorf("replies not found")
	}

	return replies, nil
}

func (repo *ReplieRepo) GetReplie(replieId int) (*model.Replie, error) {
	replie := &model.Replie{}

	rows, err := repo.db.Query("SELECT * FROM replie WHERE id=$1;", replieId)
	if err != nil {
		return nil, fmt.Errorf("error while getting replie: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&replie.ID, &replie.UserID, &replie.MessageID, &replie.Title)
		if err != nil {
			continue
		}
	}

	if replie.Title == "" {
		return nil, fmt.Errorf("replie not found")
	}

	return replie, nil
}

func (repo *ReplieRepo) GetReplieByName(name string) (*model.Replie, error) {
	replie := &model.Replie{}

	rows, err := repo.db.Query("SELECT * FROM replie WHERE title=$1;", name)
	if err != nil {
		return nil, fmt.Errorf("error while getting replie: %w", err)
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
	row := repo.db.QueryRow("INSERT INTO replie (user_id, message_id, title) VALUES ($1, $2, $3) RETURNING id;", replie.UserID, replie.MessageID, replie.Title)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while creating replie: %w", err)
	}

	return 1, nil
}

func (repo *ReplieRepo) DeleteReplie(replieId int) error {
	_, err := repo.db.Exec("DELETE FROM replie WHERE id=$1;", replieId)
	if err != nil {
		return fmt.Errorf("error while deleting replie: %w", err)
	}

	return nil
}
