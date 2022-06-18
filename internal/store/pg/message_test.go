package pg_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store/pg"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestMessagePg_CreateMessage(t *testing.T) {
	db, mock, err := utils.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   model.Message
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO message(channel_id, user_id, title, message_url, image) VALUES ($1, $2, $3, $4, $5) RETURNING id;").
					WithArgs(1, 1, "test", "test.com", "test.jpg").WillReturnRows(rows)
			},
			input: model.Message{ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test.com", Image: "test.jpg"},
			want:  1,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO message(channel_id, user_id, title, message_url, image) VALUES ($1, $2, $3, $4, $5) RETURNING id;").
					WithArgs().WillReturnRows(rows)
			},
			input:   model.Message{},
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

func TestMessagePg_GetMessageByName(t *testing.T) {
	db, mock, err := utils.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.Message
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"}).
					AddRow(1, 1, 1, "test1")

				mock.ExpectQuery("SELECT * FROM message WHERE title=$1;").
					WithArgs("test1").WillReturnRows(rows)
			},
			input: "test1",
			want:  &model.Message{ID: 1, UserID: 1, ChannelID: 1, Title: "test1"},
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"})

				mock.ExpectQuery("SELECT * FROM message WHERE title=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
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
	db, mock, err := utils.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		want    []model.Message
		wantErr bool
	}{
		{
			name: "Ok: [Replies found]",
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
			name: "Error: [Replies not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "count"})

				mock.ExpectQuery("SELECT m.id, COUNT(r.id) FROM message m LEFT JOIN replie r ON r.message_id = m.id GROUP BY m.id ORDER BY m.id;").
					WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "Error: [PQ error]",
			mock: func() {
				mock.ExpectQuery("SELECT m.id, COUNT(r.id) FROM message m LEFT JOIN replie r ON r.message_id = m.id GROUP BY m.id ORDER BY m.id;").
					WillReturnError(sqlmock.ErrCancelled)
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
