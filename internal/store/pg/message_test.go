package pg

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMessagePg_CreateMessage(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

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

				mock.ExpectQuery("INSERT INTO message(channel_id,title) VALUES ($1,$2) RETURNING id;").
					WithArgs(1, "test").WillReturnRows(rows)
			},
			input: model.Message{ChannelID: 1, Title: "test"},
			want:  1,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO message(channel_id,title) VALUES ($1,$2) RETURNING id;").
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

func TestMessagePg_GetMessage(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    *model.Message
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "channel_id", "title"}).
					AddRow(1, 1, "test")

				mock.ExpectQuery("SELECT * FROM message WHERE id=$1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.Message{ID: 1, ChannelID: 1, Title: "test"},
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "channel_id", "title"})

				mock.ExpectQuery("SELECT * FROM message WHERE id=$1;").
					WithArgs(404).WillReturnRows(rows)
			},
			input:   404,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessage(tt.input)
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

func TestMessagePg_GetMessages(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		want    []model.Message
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "channel_id", "title"}).
					AddRow(1, 1, "test1").
					AddRow(2, 2, "test2")

				mock.ExpectQuery("SELECT * FROM message;").
					WillReturnRows(rows)
			},
			want: []model.Message{
				{ID: 1, ChannelID: 1, Title: "test1"},
				{ID: 2, ChannelID: 2, Title: "test2"},
			},
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "channel_id", "title"})

				mock.ExpectQuery("SELECT * FROM message;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessages()
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
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

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
				rows := sqlmock.NewRows([]string{"id", "channel_id", "title"}).
					AddRow(1, 1, "test1")

				mock.ExpectQuery("SELECT * FROM message WHERE title=$1;").
					WithArgs("test1").WillReturnRows(rows)
			},
			input: "test1",
			want:  &model.Message{ID: 1, ChannelID: 1, Title: "test1"},
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "channel_id", "title"})

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

func TestMessagePg_DeleteMessage(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM message WHERE id=$1;").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: 1,
		},
		{
			name: "message not found",
			mock: func() {
				mock.ExpectExec("DELETE FROM message WHERE id=$1;").
					WithArgs(404).WillReturnError(sql.ErrNoRows)
			},
			input:   404,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteMessage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
