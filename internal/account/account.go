package account

import (
	// Core packages
	"context"
	"database/sql"
	"time"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/auth"

	// Third-party packages
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrAuthenticationFailure occurs when an account attempts to authenticate but something goes wrong.
	ErrAuthenticationFailure = errors.New("Authentication failed")
)

/**
 * Authenticate finds an account by the name (unique) and verifies the password. On
 * success it returns a Claims value representing the account. The claims can be
 * used to generate a token for future authentication.
 */

func Authenticate(ctx context.Context, db *sqlx.DB, now time.Time, name, password string) (auth.Claims, error) {

	const q = `SELECT * FROM account WHERE name = $1`

	var a Account
	if err := db.GetContext(ctx, &a, q, name); err != nil {

		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated account which emails are in the system.
		if err == sql.ErrNoRows {
			return auth.Claims{}, ErrAuthenticationFailure
		}

		return auth.Claims{}, errors.Wrap(err, "selecting single account")
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(a.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the account
	// and generate the token for the account.
	claims := auth.NewClaims(a.Id, a.Roles, now, time.Hour)
	return claims, nil
}

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
