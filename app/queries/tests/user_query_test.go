package queries_test

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/app/queries"
	"github.com/create-go-app/fiber-go-template/pkg/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func TestGetUserByID(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	q := &queries.UserQueries{DB: db}

	// Sample data
	id := uuid.New()
	user := models.User{
		ID:           id,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		UserStatus:   1,
		UserRole:     repository.AdminRoleName,
	}

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "email", "password_hash", "user_status", "user_role",
	}).AddRow(
		user.ID, user.CreatedAt, user.UpdatedAt, user.Email, user.PasswordHash, user.UserStatus, user.UserRole,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)

	got, err := q.GetUserByID(id)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Email, got.Email)

	// Test no rows
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE id = $1`)).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, err = q.GetUserByID(id)
	assert.Error(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	q := &queries.UserQueries{DB: db}

	email := "test@example.com"
	user := models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        email,
		PasswordHash: "hashed",
		UserStatus:   1,
		UserRole:     repository.AdminRoleName,
	}

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "email", "password_hash", "user_status", "user_role",
	}).AddRow(
		user.ID, user.CreatedAt, user.UpdatedAt, user.Email, user.PasswordHash, user.UserStatus, user.UserRole,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE email = $1`)).
		WithArgs(email).
		WillReturnRows(rows)

	got, err := q.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, got.Email)

	// Test no rows
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE email = $1`)).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	_, err = q.GetUserByEmail(email)
	assert.Error(t, err)
}

func TestCreateUser(t *testing.T) {
	db, mock := newMockDB(t)
	defer db.Close()

	q := &queries.UserQueries{DB: db}

	user := &models.User{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        "test@example.com",
		PasswordHash: "hashed",
		UserStatus:   1,
		UserRole:     repository.AdminRoleName,
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7)`)).
		WithArgs(user.ID, user.CreatedAt, user.UpdatedAt, user.Email, user.PasswordHash, user.UserStatus, user.UserRole).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := q.CreateUser(user)
	assert.NoError(t, err)

	// Test Exec error
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7)`)).
		WithArgs(user.ID, user.CreatedAt, user.UpdatedAt, user.Email, user.PasswordHash, user.UserStatus, user.UserRole).
		WillReturnError(assert.AnError)

	err = q.CreateUser(user)
	assert.Error(t, err)
}
