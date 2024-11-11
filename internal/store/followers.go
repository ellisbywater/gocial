package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Followers struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) Follow(ctx context.Context, followerId int64, userId int64) error {
	query := `
		INSERT INTO followers(
			user_id,
			follower_id
		) Values ($1, $2);
	`
	ctx, cancel := context.WithTimeout(ctx, QueryCtxTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}
	return err
}

func (s *FollowersStore) Unfollow(ctx context.Context, followerId int64, userId int64) error {
	query := `
		DELETE followers WHERE user_id=$1 AND follower_id = $2;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryCtxTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	return err
}
