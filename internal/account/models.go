package account

import (
	"time"

	"github.com/lib/pq"
)

// Account represents a thing with access to the system.
type Account struct {
	Id           string         `db:"id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Roles        pq.StringArray `db:"roles" json:"roles"`
	PasswordHash []byte         `db:"password_hash" json:"-"`
	DateCreated  time.Time      `db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `db:"date_updated" json:"date_updated"`
}

// NewAccount contains information needed to create a new Account.
type NewAccount struct {
	Name            string   `json:"name" validate:"required"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"password_confirm" validate:"eqfield=Password"`
}
