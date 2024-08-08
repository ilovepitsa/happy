package repo

import (
	"context"

	"github.com/ilovepitsa/happy/notify/pkg/config"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	NotifRepo *NotificationRepo
}

func NewRepo(cfg config.Postgres) (*Repo, error) {

	conntection, err := pgx.Connect(context.TODO(), cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	nRepo := NewNotificationRepo(conntection)

	return &Repo{
		NotifRepo: nRepo,
	}, nil
}
