package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ellisbywater/gocial/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "users"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	//
	follower := getUserFromContext(r)

	// revert at auth implementation
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.store.Followers.Follow(ctx, follower.ID, payload.UserID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictingResourceError(w, r, err)
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, follower); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	follower := getUserFromContext(r)

	//TODO: revert at auth implementation
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.store.Followers.Unfollow(ctx, follower.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, follower); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.badRequestError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
