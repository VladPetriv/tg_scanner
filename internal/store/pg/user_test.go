package pg_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store/pg"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestUserPg_CreateUser(t *testing.T) {
	db, mock, err := utils.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   model.User
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO user (username, firstname, lastname, photourl) VALUES ($1, $2, $3, $4) RETURNING id;").
					WithArgs("test", "test", "test", "test.jpg").WillReturnRows(rows)
			},
			input: model.User{Username: "test", FirstName: "test", LastName: "test", PhotoURL: "test.jpg"},
			want:  1,
		},
		{
			name: "Empty fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO user (username, firstname, lastname, photourl) VALUES ($1, $2, $3, $4) RETURNING id;").
					WithArgs("", "", "", "").WillReturnRows(rows)
			},
			input:   model.User{},
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

func TestUserPg_GetUsers(t *testing.T) {
	db, mock, err := utils.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		want    []model.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "firstname", "lastname", "photourl"}).
					AddRow(1, "test1", "test1", "test1", "test1.jpg").
					AddRow(2, "test2", "test2", "test2", "test2.jpg")

				mock.ExpectQuery("SELECT * FROM user;").
					WillReturnRows(rows)
			},
			want: []model.User{
				{ID: 1, Username: "test1", FirstName: "test1", LastName: "test1", PhotoURL: "test1.jpg"},
				{ID: 2, Username: "test2", FirstName: "test2", LastName: "test2", PhotoURL: "test2.jpg"},
			},
		},
		{
			name: "users not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "firstname", "lastname", "photourl"})

				mock.ExpectQuery("SELECT * FROM user;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUsers()

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

func TestUserPg_GetUserByUsername(t *testing.T) {
	db, mock, err := utils.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "firstname", "lastname", "photourl"}).
					AddRow(1, "test", "test", "test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM user WHERE username=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want: &model.User{
				ID:        1,
				Username:  "test",
				FirstName: "test",
				LastName:  "test",
				PhotoURL:  "test.jpg",
			},
		},
		{
			name: "Empty fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "firstname", "lastname", "photourl"})

				mock.ExpectQuery("SELECT * FROM user WHERE username=$1;").
					WithArgs().WillReturnRows(rows)
			},
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
