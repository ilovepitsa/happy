package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrFollowUserNotExist = errors.New("one of user doesnt exist")
	ErrAddFollow          = errors.New("cant add follow")
)

type FollowRepo struct {
	connection *pgx.Conn
}

func NewFollowRepo(connection *pgx.Conn) *FollowRepo {
	return &FollowRepo{
		connection: connection,
	}
}

func (fr *FollowRepo) AddFollow(userId int64, target int64, notifyBefore string) error {
	trans, err := fr.connection.Begin(context.TODO())
	if err != nil {
		return err
	}

	res, err := trans.Query(context.TODO(), "Insert Into follows (subscriber_id, target_id, notify_before) values ($1, $2, $3) RETURNING 1;", userId, target, notifyBefore)
	if err != nil {
		return err
	}

	if res.CommandTag().RowsAffected() == 0 {
		return ErrAddFollow
	}
	return nil
}

func (fr *FollowRepo) RemoveFollow()
