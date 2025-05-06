package middlewares

import (
	"context"
	"net/http"
	"task2/helperFunc"

	"github.com/go-playground/validator/v10"
)

type contextKey string

const BindKey contextKey = "parsedBody"

var validate = validator.New()

func BindAndValidate(schema interface{}, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := helperFunc.Clone(schema)

		if err := helperFunc.FromJSON(reqBody, r.Body); err != nil {
			helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid JSON", nil)
			return
		}

		if err := validate.Struct(reqBody); err != nil {
			helperFunc.RespondJSON(w, http.StatusBadRequest, "Validation error", err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), BindKey, reqBody)
		next(w, r.WithContext(ctx))
	}
}
