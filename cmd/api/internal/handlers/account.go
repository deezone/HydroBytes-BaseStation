package handlers

import (
	// Core packages
	"context"
	"net/http"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/auth"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/web"
	"github.com/deezone/HydroBytes-BaseStation/internal/account"

	// Third-party packages
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Account holds handlers for dealing with an account.
type Account struct {
	db            *sqlx.DB
	authenticator *auth.Authenticator
}

// Token generates an authentication token for an account. The client must include
// a name and password for the request using HTTP Basic Auth. The account will
// be identified by name and authenticated by the password.
func (a *Account) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	ctx, span := trace.StartSpan(ctx, "handlers.Account.Token")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return errors.New("web value missing from context")
	}

	name, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide name and password in Basic auth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := account.Authenticate(ctx, a.db, v.Start, name, pass)
	if err != nil {
		switch err {
		case account.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "authenticating")
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = a.authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
