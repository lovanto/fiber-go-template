package queries_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/app/queries"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func TestBookQueries_GetBooks(t *testing.T) {
	db, mock := newMock(t)
	q := &queries.BookQueries{DB: db}

	columns := []string{"id", "created_at", "updated_at", "user_id", "title", "author", "book_status", "book_attrs"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM books`)).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(uuid.New(), time.Now(), time.Now(), uuid.New(), "Title1", "Author1", 1, []byte(`{"picture": "picture1", "description": "description1", "rating": 5}`)))

	books, err := q.GetBooks()
	assert.NoError(t, err)
	assert.Len(t, books, 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM books`)).
		WillReturnError(errors.New("db error"))
	_, err = q.GetBooks()
	assert.Error(t, err)
}

func TestBookQueries_GetBooksByAuthor(t *testing.T) {
	db, mock := newMock(t)
	q := &queries.BookQueries{DB: db}

	columns := []string{"id", "created_at", "updated_at", "user_id", "title", "author", "book_status", "book_attrs"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM books WHERE author = $1`)).
		WithArgs("Author1").
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(uuid.New(), time.Now(), time.Now(), uuid.New(), "Title1", "Author1", 1, []byte(`{"picture": "picture1", "description": "description1", "rating": 5}`)))

	books, err := q.GetBooksByAuthor("Author1")
	assert.NoError(t, err)
	assert.Len(t, books, 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM books WHERE author = $1`)).
		WithArgs("Author1").
		WillReturnError(errors.New("db error"))
	_, err = q.GetBooksByAuthor("Author1")
	assert.Error(t, err)
}

func TestBookQueries_GetBook(t *testing.T) {
	db, mock := newMock(t)
	q := &queries.BookQueries{DB: db}
	id := uuid.New()

	columns := []string{"id", "created_at", "updated_at", "user_id", "title", "author", "book_status", "book_attrs"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM books WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(id, time.Now(), time.Now(), uuid.New(), "Title1", "Author1", 1, []byte(`{"picture": "picture1", "description": "description1", "rating": 5}`)))

	book, err := q.GetBook(id)
	assert.NoError(t, err)
	assert.Equal(t, id, book.ID)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM books WHERE id = $1`)).
		WithArgs(id).
		WillReturnError(errors.New("db error"))
	_, err = q.GetBook(id)
	assert.Error(t, err)
}

func TestBookQueries_CreateBook(t *testing.T) {
	db, mock := newMock(t)
	q := &queries.BookQueries{DB: db}

	b := &models.Book{
		ID:         uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		UserID:     uuid.New(),
		Title:      "Title1",
		Author:     "Author1",
		BookStatus: 1,
		BookAttrs: models.BookAttrs{
			Picture:     "Picture1",
			Description: "Description1",
			Rating:      5,
		},
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO books VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)).
		WithArgs(b.ID, b.CreatedAt, b.UpdatedAt, b.UserID, b.Title, b.Author, b.BookStatus, b.BookAttrs).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := q.CreateBook(b)
	assert.NoError(t, err)

	// error case
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO books VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)).
		WithArgs(b.ID, b.CreatedAt, b.UpdatedAt, b.UserID, b.Title, b.Author, b.BookStatus, b.BookAttrs).
		WillReturnError(errors.New("insert error"))
	err = q.CreateBook(b)
	assert.Error(t, err)
}

func TestBookQueries_UpdateBook(t *testing.T) {
	db, mock := newMock(t)
	q := &queries.BookQueries{DB: db}
	id := uuid.New()
	b := &models.Book{
		UpdatedAt:  time.Now(),
		Title:      "UpdatedTitle",
		Author:     "UpdatedAuthor",
		BookStatus: 1,
		BookAttrs: models.BookAttrs{
			Picture:     "UpdatedPicture",
			Description: "UpdatedDescription",
			Rating:      5,
		},
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET updated_at = $2, title = $3, author = $4, book_status = $5, book_attrs = $6 WHERE id = $1`)).
		WithArgs(id, b.UpdatedAt, b.Title, b.Author, b.BookStatus, b.BookAttrs).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := q.UpdateBook(id, b)
	assert.NoError(t, err)

	// error case
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE books SET updated_at = $2, title = $3, author = $4, book_status = $5, book_attrs = $6 WHERE id = $1`)).
		WithArgs(id, b.UpdatedAt, b.Title, b.Author, b.BookStatus, b.BookAttrs).
		WillReturnError(errors.New("update error"))
	err = q.UpdateBook(id, b)
	assert.Error(t, err)
}

func TestBookQueries_DeleteBook(t *testing.T) {
	db, mock := newMock(t)
	q := &queries.BookQueries{DB: db}
	id := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := q.DeleteBook(id)
	assert.NoError(t, err)

	// error case
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM books WHERE id = $1`)).
		WithArgs(id).
		WillReturnError(errors.New("delete error"))
	err = q.DeleteBook(id)
	assert.Error(t, err)
}
