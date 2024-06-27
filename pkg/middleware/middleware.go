package middleware

import (
	"net/http"

	"github.com/http-everything/httpe/pkg/auth"
	"github.com/http-everything/httpe/pkg/response"
	"github.com/http-everything/httpe/pkg/rules"
	"github.com/http-everything/httpe/pkg/share/firstof"
	"github.com/http-everything/httpe/pkg/share/logger"

	"github.com/dustin/go-humanize" //nolint:misspell
)

const DefaultMaxRequestBody = "512KB"

type Middleware struct {
	rule   rules.Rule
	logger *logger.Logger
}

func New(rule rules.Rule, logger *logger.Logger) Middleware {
	return Middleware{
		rule:   rule,
		logger: logger,
	}
}

func (m Middleware) Collection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Initialise a new http response writer.
		respWriter := response.New(w, m.rule.Respond, m.logger)

		// Reject requests exceeding the max request body limit
		lim, err := humanize.ParseBytes(firstof.String(m.rule.MaxRequestBody(), DefaultMaxRequestBody))
		if err != nil {
			respWriter.InternalServerErrorf(
				"error parsing max_request_body '%s': %s",
				m.rule.With.MaxRequestBody,
				err.Error(),
			)
			return
		}
		if int(r.ContentLength) > int(lim) {
			respWriter.RequestEntityTooLarge(int(r.ContentLength), int(lim))
			return
		}

		// Authenticate, if requested by the rule
		if m.rule.With != nil {
			ok, err := auth.IsRequestAuthenticated(m.rule.With.AuthBasic, m.rule.With.AuthHashing, r)
			if err != nil {
				respWriter.InternalServerError(err)
				return
			}
			if !ok {
				respWriter.Unauthorised()
				return
			}
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
