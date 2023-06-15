package middlewares

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gotodo/config"
	"github.com/gotodo/helpers"
	"github.com/gotodo/models"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				helpers.ResponseJSON(w, models.ResponseBody{Message: "Unauthorized", Code: http.StatusUnauthorized})
				return
			}
		}

		tokenString := c.Value

		claims := &config.JWTClaim{}

		//Parsing token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorExpired:
				// token expired
				helpers.ResponseJSON(w, models.ResponseBody{Message: "Unauthorized. Token Expired", Code: http.StatusUnauthorized})
				return
			default:
				helpers.ResponseJSON(w, models.ResponseBody{Message: "Unauthorized", Code: http.StatusUnauthorized})
				return
			}
		}

		if !token.Valid {
			helpers.ResponseJSON(w, models.ResponseBody{Message: "Unauthorized", Code: http.StatusUnauthorized})
			return
		}

		ctx := context.WithValue(r.Context(), config.ContextUserKey, claims.Email)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}
