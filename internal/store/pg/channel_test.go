package pg_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store/pg"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestChannelPg_CreateChannel(t *testing.T) {
	db, mock, err := utils.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewChannelRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   model.Channel
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO channel(name, title, photourl) VALUES ($1, $2, $3) RETURNING id;").
					WithArgs("test", "test", "test.jpg").WillReturnRows(rows)

			},
			input: model.Channel{ID: 0, Name: "test", Title: "test", PhotoURL: "test.jpg"},
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateChannel(&tt.input)
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

func TestChannelPg_GetChannelByName(t *testing.T) {
	db, mock, err := utils.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewChannelRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.Channel
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"}).
					AddRow(1, "test", "test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM channel WHERE name=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.Channel{ID: 1, Name: "test", Title: "test", PhotoURL: "test.jpg"},
		},
		{
			name: "channel not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"})

				mock.ExpectQuery("SELECT * FROM channel WHERE name=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelByName(tt.input)
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
