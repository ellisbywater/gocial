package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ellisbywater/gocial/internal/mailer"
	"github.com/ellisbywater/gocial/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// store user
	ctx := r.Context()
	pToken := uuid.New().String()
	hash := sha256.Sum256([]byte(pToken))
	hashToken := hex.EncodeToString(hash[:])
	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mailer.exp)

	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: pToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, pToken)
	isProdEnv := app.config.env == "production"

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// mail
	status, err := app.mailer.Send(mailer.UserEmailActivationTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending activation email", "error", err)
		// rollback user creation if email fails
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}
		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent.", "status_code", status)
	log.Printf("Email sent with status code")

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createUserTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedError(w, r, err)
		return
	}

	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}

}
