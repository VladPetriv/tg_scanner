package pg_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store/pg"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserPg_GetUserByUsername(t *testing.T) {
	dbM, mock, err := utils.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer dbM.Close()

	db := sqlx.NewDb(dbM, "postgres")

	r := pg.NewUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.User
		wantErr bool
	}{
		{
			name: "Ok: [user found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"}).
					AddRow(1, "test", "test test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM tg_user WHERE username=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.User{ID: 1, Username: "test", FullName: "test test", ImageURL: "test.jpg"},
		},
		{
			name: "Error: [user not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"})

				mock.ExpectQuery("SELECT * FROM tg_user WHERE username=$1;").
					WithArgs().WillReturnRows(rows)
			},
			input: "test",
			want:  nil,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM tg_user WHERE username=$1;").
					WithArgs().WillReturnError(fmt.Errorf("some error"))
			},
			input:   "test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserByUsername(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPg_CreateUser(t *testing.T) {
	dbM, mock, err := utils.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer dbM.Close()

	db := sqlx.NewDb(dbM, "postgres")

	r := pg.NewUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   model.User
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [user created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO tg_user (username, fullname, imageurl) VALUES ($1, $2, $3) RETURNING id;").
					WithArgs("test", "test test", "test.jpg").WillReturnRows(rows)
			},
			input: model.User{Username: "test", FullName: "test test", ImageURL: "test.jpg"},
			want:  1,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("INSERT INTO tg_user (username, fullname, imageurl) VALUES ($1, $2, $3) RETURNING id;").
					WithArgs("test", "test test", "test.jpg").WillReturnError(fmt.Errorf("some error"))
			},
			input:   model.User{Username: "test", FullName: "test test", ImageURL: "test.jpg"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateUser(&tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
