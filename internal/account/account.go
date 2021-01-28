package account

import (
	// Core packages
	"context"
	"time"

	// Third-party packages
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Create inserts a new account into the database.
func Create(ctx context.Context, db *sqlx.DB, n NewAccount, now time.Time) (*Account, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	a := Account{
		Id:           uuid.New().String(),
		Name:         n.Name,
		PasswordHash: hash,
		Roles:        n.Roles,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `INSERT INTO account
		(id, name, password_hash, roles, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.ExecContext(
		ctx, q,
		a.Id, a.Name,
		a.PasswordHash, a.Roles,
		a.DateCreated, a.DateUpdated,
	)
	if err != nil {
		return nil, errors.Wrap(err, "inserting account")
	}

	return &a, nil
}
