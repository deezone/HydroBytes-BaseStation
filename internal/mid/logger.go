package mid

import (
	// Core packages
	"errors"
	"log"
	"net/http"
	"time"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/web"
)

// Logger writes some information about the request to the logs in the
// format: (200) GET /foo -> IP ADDR (latency)
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {
			v, ok := r.Context().Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web value missing from context")
			}

			err := before(w, r)

			// Format log message
			// ex: POST (201) : /v1/station-type -> 127.0.0.1:53670 (6.408225ms)
			log.Printf("%s (%d) : %s -> %s (%s)",
				r.Method, v.StatusCode,
				r.URL.Path,
				r.RemoteAddr, time.Since(v.Start),
			)

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return f
}
