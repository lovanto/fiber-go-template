package queries

import (
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserQueries struct {
	*sqlx.DB
}

func (q *UserQueries) GetUserByID(id uuid.UUID) (models.User, error) {
	user := models.User{}
	query := `SELECT * FROM users WHERE id = $1`

	err := q.Get(&user, query, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}
	query := `SELECT * FROM users WHERE email = $1`

	err := q.Get(&user, query, email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserQueries) CreateUser(u *models.User) error {
	query := `INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := q.Exec(
		query,
		u.ID, u.CreatedAt, u.UpdatedAt, u.Email, u.PasswordHash, u.UserStatus, u.UserRole,
	)
	if err != nil {
		return err
	}

	return nil
}
