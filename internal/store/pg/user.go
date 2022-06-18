package pg

import (
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
	user := &model.User{}

	rows, err := repo.db.Query("SELECT * FROM tg_user WHERE username=$1;", username)
	if err != nil {
		return nil, &utils.GettingError{Name: "user by username", ErrorValue: err}
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.PhotoURL)
		if err != nil {
			continue
		}
	}

	if user.Username == "" && user.FullName == "" {
		return nil, nil
	}

	return user, nil
}

func (repo *UserRepo) CreateUser(user *model.User) (int, error) {
	var id int

	row := repo.db.QueryRow(
		"INSERT INTO tg_user (username, fullname, photourl) VALUES ($1, $2, $3) RETURNING id;",
		user.Username, user.FullName, user.PhotoURL,
	)
	if err := row.Scan(&id); err != nil {
		return id, &utils.CreateError{Name: "user", ErrorValue: err}
	}

	return id, nil
}
