package pg

import (
	"database/sql"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) GetUserByUsername(username string) (*model.User, error) {
	user := model.User{}

	err := repo.db.Get(&user, "SELECT * FROM tg_user WHERE username=$1;", username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, &utils.GettingError{Name: "user by username", ErrorValue: err}
	}

	return &user, nil
}

func (repo *UserRepo) CreateUser(user *model.User) (int, error) {
	var id int

	row := repo.db.QueryRow(
		"INSERT INTO tg_user (username, fullname, imageurl) VALUES ($1, $2, $3) RETURNING id;",
		user.Username, user.FullName, user.ImageURL,
	)
	if err := row.Scan(&id); err != nil {
		return id, &utils.CreateError{Name: "user", ErrorValue: err}
	}

	return id, nil
}
