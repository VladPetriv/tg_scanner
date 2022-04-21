package pg

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) GetUsers() ([]model.User, error) {
	users := make([]model.User, 0)
	rows, err := repo.db.Query("SELECT * FROM tg_user;")
	if err != nil {
		return nil, fmt.Errorf("error while getting users: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.PhotoURL)
		if err != nil {
			continue
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("users not found")
	}

	return users, nil
}

func (repo *UserRepo) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}

	rows, err := repo.db.Query("SELECT * FROM tg_user WHERE username=$1;", username)
	if err != nil {
		return nil, fmt.Errorf("error while getting user by username: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.PhotoURL)
		if err != nil {
			continue
		}
	}

	if user.Username == "" || user.FullName == "" {
		return nil, nil
	}

	return user, nil
}

func (repo *UserRepo) CreateUser(user *model.User) (int, error) {
	var id int
	row := repo.db.QueryRow("INSERT INTO tg_user (username, fullname, photourl) VALUES ($1, $2, $3) RETURNING id;", user.Username, user.FullName, user.PhotoURL)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while creating user: %w", err)
	}

	return id, nil
}
