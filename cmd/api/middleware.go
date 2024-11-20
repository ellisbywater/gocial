package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}
		parts := strings.Split(authHeader, " ") // authorization: Bearer <token>
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is incorrect"))
			return
		}
		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userId)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}
			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is incorrect"))
				return
			}
			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			username := app.config.auth.basic.username
			pass := app.config.auth.basic.pass
			creds := strings.SplitN(string(decoded), ":", 2)
			// check the credentials
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedError(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
