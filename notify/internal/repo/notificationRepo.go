package repo

import "github.com/jackc/pgx/v5"

type NotificationRepo struct {
	connection *pgx.Conn
}

func NewNotificationRepo(connection *pgx.Conn) *NotificationRepo {
	return &NotificationRepo{
		connection: connection,
	}
}
