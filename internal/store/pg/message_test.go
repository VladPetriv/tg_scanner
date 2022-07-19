package pg_test

import (
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store/pg"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMessagePg_GetMessageByName(t *testing.T) {
	dbM, mock, err := utils.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer dbM.Close()

	db := sqlx.NewDb(dbM, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.Message
		wantErr bool
	}{
		{
			name: "Ok: [message found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "channel_id", "user_id", "title", "message_url", "imageurl"}).
					AddRow(1, 1, 1, "test1", "test.com", "test.jpg")

				mock.ExpectQuery("SELECT * FROM message WHERE title=$1;").
					WithArgs("test1").WillReturnRows(rows)
			},
			input: "test1",
			want:  &model.Message{ID: 1, UserID: 1, ChannelID: 1, Title: "test1", MessageURL: "test.com", ImageURL: "test.jpg"},
		},
		{
			name: "Error: [message not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "channel_id", "user_id", "title", "message_url", "imageurl"})

				mock.ExpectQuery("SELECT * FROM message WHERE title=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM message WHERE title=$1;").
					WithArgs().WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessageByName(tt.input)
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

func TestMessagePg_GetMessagesWithReplies(t *testing.T) {
	dbM, mock, err := utils.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer dbM.Close()

	db := sqlx.NewDb(dbM, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		want    []model.Message
		wantErr bool
	}{
		{
			name: "Ok: [replies found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "count"}).
					AddRow(1, 5).
					AddRow(2, 1)

				mock.ExpectQuery("SELECT m.id, COUNT(r.id) FROM message m LEFT JOIN replie r ON r.message_id = m.id GROUP BY m.id ORDER BY m.id;").
					WillReturnRows(rows)
			},
			want: []model.Message{
				{ID: 1, RepliesCount: 5},
				{ID: 2, RepliesCount: 1},
			},
		},
		{
			name: "Error: [replies not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "count"})

				mock.ExpectQuery("SELECT m.id, COUNT(r.id) FROM message m LEFT JOIN replie r ON r.message_id = m.id GROUP BY m.id ORDER BY m.id;").
					WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT m.id, COUNT(r.id) FROM message m LEFT JOIN replie r ON r.message_id = m.id GROUP BY m.id ORDER BY m.id;").
					WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessagesWithRepliesCount()
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

func TestMessagePg_CreateMessage(t *testing.T) {
	dbM, mock, err := utils.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer dbM.Close()

	db := sqlx.NewDb(dbM, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   model.Message
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [message created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO message(channel_id, user_id, title, message_url, imageurl) VALUES ($1, $2, $3, $4, $5) RETURNING id;").
					WithArgs(1, 1, "test", "test.com", "test.jpg").WillReturnRows(rows)
			},
			input: model.Message{ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test.com", ImageURL: "test.jpg"},
			want:  1,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("INSERT INTO message(channel_id, user_id, title, message_url, imageurl) VALUES ($1, $2, $3, $4, $5) RETURNING id;").
					WithArgs().WillReturnError(fmt.Errorf("some error"))
			},
			input:   model.Message{ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test.com", ImageURL: "test.jpg"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateMessage(&tt.input)
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
