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

				mock.ExpectQuery("INSERT INTO tg_user (username, fullname, photourl) VALUES ($1, $2, $3) RETURNING id;").
					WithArgs("test", "test test", "test.jpg").WillReturnRows(rows)
			},
			input: model.User{Username: "test", FullName: "test test", PhotoURL: "test.jpg"},
			want:  1,
		},
		{
			name: "Empty fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO tg_user (username, fullname, photourl) VALUES ($1, $2, $3) RETURNING id;").
					WithArgs("", "", "").WillReturnRows(rows)
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
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "photourl"}).
					AddRow(1, "test1", "test test1", "test1.jpg").
					AddRow(2, "test2", "test test2", "test2.jpg")

				mock.ExpectQuery("SELECT * FROM tg_user;").
					WillReturnRows(rows)
			},
			want: []model.User{
				{ID: 1, Username: "test1", FullName: "test test1", PhotoURL: "test1.jpg"},
				{ID: 2, Username: "test2", FullName: "test test2", PhotoURL: "test2.jpg"},
			},
		},
		{
			name: "users not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "photourl"})

				mock.ExpectQuery("SELECT * FROM tg_user;").
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
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "photourl"}).
					AddRow(1, "test", "test test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM tg_user WHERE username=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want: &model.User{
				ID:       1,
				Username: "test",
				FullName: "test test",
				PhotoURL: "test.jpg",
			},
		},
		{
			name: "Empty fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "photourl"})

				mock.ExpectQuery("SELECT * FROM tg_user WHERE username=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
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
