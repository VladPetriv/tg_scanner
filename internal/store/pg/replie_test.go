package pg

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRepliePg_CreateReplie(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewReplieRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   model.Replie
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO replie (message_id, title) VALUES ($1, $2) RETURNING id;").
					WithArgs(1, "test").WillReturnRows(rows)

			},
			input: model.Replie{MessageID: 1, Title: "test"},
			want:  1,
		},
		{
			name: "empty fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO replie (message_id, title) VALUES ($1, $2) RETURNING id;").
					WithArgs().WillReturnRows(rows)
			},
			input:   model.Replie{},
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

func TestRepliePg_GetReplie(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewReplieRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    *model.Replie
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message_id", "title"}).
					AddRow(1, 1, "test")

				mock.ExpectQuery("SELECT * FROM replie WHERE id=$1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.Replie{ID: 1, MessageID: 1, Title: "test"},
		},
		{
			name: "replie not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message_id", "title"})

				mock.ExpectQuery("SELECT * FROM replie WHERE id=$1;").
					WithArgs().WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetReplie(tt.input)
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

func TestRepliePg_GetReplies(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewReplieRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		want    []model.Replie
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message_id", "title"}).
					AddRow(1, 1, "test1").
					AddRow(2, 2, "test2")

				mock.ExpectQuery("SELECT * FROM replie;").
					WillReturnRows(rows)
			},
			want: []model.Replie{
				{ID: 1, MessageID: 1, Title: "test1"},
				{ID: 2, MessageID: 2, Title: "test2"},
			},
		},
		{
			name: "replies not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message_id", "title"})

				mock.ExpectQuery("SELECT * FROM replie;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetReplies()
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

func TestRepliePg_GetReplieByName(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewReplieRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.Replie
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message_id", "title"}).
					AddRow(1, 1, "test")

				mock.ExpectQuery("SELECT * FROM replie WHERE title=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.Replie{ID: 1, MessageID: 1, Title: "test"},
		},
		{
			name: "replie not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message_id", "title"})

				mock.ExpectQuery("SELECT * FROM replie WHERE title=$1;").
					WithArgs().WillReturnRows(rows)
			},
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

func TestRepliePg_DeleteReplie(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewReplieRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM replie WHERE id=$1;").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: 1,
		},
		{
			name: "replie not found",
			mock: func() {
				mock.ExpectExec("DELETE FROM replie WHERE id=$1;").
					WithArgs(404).WillReturnError(sql.ErrNoRows)
			},
			input:   404,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteReplie(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
