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

func TestRepliePg_GetReplieByName(t *testing.T) {
	dbM, mock, err := utils.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer dbM.Close()

	db := sqlx.NewDb(dbM, "postgres")

	r := pg.NewReplieRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.Replie
		wantErr bool
	}{
		{
			name: "Ok: [replie found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id", "title", "imageurl"}).
					AddRow(1, 1, 1, "test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM replie WHERE title=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.Replie{ID: 1, UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"},
		},
		{
			name: "Error: [replie not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id", "title", "imageurl"})

				mock.ExpectQuery("SELECT * FROM replie WHERE title=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM replie WHERE title=$1;").
					WithArgs().WillReturnError(fmt.Errorf("some error"))
			},
			input:   "test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetReplieByName(tt.input)
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

func TestRepliePg_CreateReplie(t *testing.T) {
	dbM, mock, err := utils.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer dbM.Close()

	db := sqlx.NewDb(dbM, "postres")

	r := pg.NewReplieRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   model.Replie
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [replie created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO replie (user_id, message_id, title, imageurl) VALUES ($1, $2, $3, $4) RETURNING id;").
					WithArgs(1, 1, "test", "test.jpg").WillReturnRows(rows)

			},
			input: model.Replie{UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"},
			want:  1,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("INSERT INTO replie (user_id, message_id, title, imageurl) VALUES ($1, $2, $3, $4) RETURNING id;").
					WithArgs().WillReturnError(fmt.Errorf("some error"))
			},
			input:   model.Replie{UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateReplie(&tt.input)
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
