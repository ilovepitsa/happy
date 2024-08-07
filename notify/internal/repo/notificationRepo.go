package repo

import (
	"context"
	"errors"
	"time"

	"github.com/ilovepitsa/happy/notify/internal/entity"
	"github.com/jackc/pgx/v5"
)

var (
	ErrAddNotification = errors.New("cant add notification")
)

type NotificationRepo struct {
	connection *pgx.Conn
}

func NewNotificationRepo(connection *pgx.Conn) *NotificationRepo {
	return &NotificationRepo{
		connection: connection,
	}
}

func (nr *NotificationRepo) GetUserNotifications(userId int32, idAbove int) ([]entity.Notification, int, error) {
	trans, err := nr.connection.Begin(context.TODO())
	if err != nil {
		trans.Rollback(context.TODO())
		return nil, 0, err
	}
	lastId := 0
	res, err := trans.Query(context.TODO(), "select id, target, text, date from notification where type = 'web' and target = '$1' and id > $2", userId, idAbove)
	if err != nil {
		trans.Rollback(context.TODO())
		return nil, 0, err
	}
	defer res.Close()
	var result []entity.Notification

	for res.Next() {
		notification := entity.Notification{}
		err = res.Scan(&lastId, &notification.Target, &notification.NotificationText, &notification.Date)
		if err != nil {
			trans.Rollback(context.TODO())
			return nil, 0, err
		}

		result = append(result, notification)
	}

	trans.Commit(context.TODO())
	return result, lastId, nil
}

func (nr *NotificationRepo) GetUnsendEmailNotifications() ([]entity.Notification, error) {
	trans, err := nr.connection.Begin(context.TODO())
	if err != nil {
		trans.Rollback(context.TODO())
		return nil, err
	}

	res, err := trans.Query(context.TODO(), "select target, text, date from notification where send = false and type = 'email' ")
	if err != nil {
		trans.Rollback(context.TODO())
		return nil, err
	}
	defer res.Close()
	var result []entity.Notification

	for res.Next() {
		notification := entity.Notification{}
		err = res.Scan(&notification.Target, &notification.NotificationText, &notification.Date)
		if err != nil {
			trans.Rollback(context.TODO())
			return nil, err
		}

		result = append(result, notification)
	}

	trans.Commit(context.TODO())
	return result, nil
}

func (nr *NotificationRepo) GetNotificationText(isEmail bool) (string, error) {
	trans, err := nr.connection.Begin(context.TODO())
	if err != nil {
		trans.Rollback(context.TODO())
		return "", err
	}

	res, err := trans.Query(context.TODO(), "select template from notification_template where isEmail = $1';", isEmail)
	if err != nil {
		trans.Rollback(context.TODO())
		return "", err
	}
	defer res.Close()
	var text string
	err = res.Scan(&text)
	if err != nil {
		trans.Rollback(context.TODO())
		return "", err
	}
	trans.Commit(context.TODO())
	return text, nil
}

func (nr *NotificationRepo) AddNotification(target int32, isEmail bool, text string, date time.Time, notify_before time.Duration) error {
	trans, err := nr.connection.Begin(context.TODO())
	if err != nil {
		trans.Rollback(context.TODO())
		return err
	}
	typeNotifications := "web"
	if isEmail {
		typeNotifications = "email"
	}
	res, err := trans.Query(context.TODO(), "insert into notification(target, type, text, date, send) values ($1, $2, $3, $4, false) RETURNING 1", target, typeNotifications, text, date)
	if err != nil {
		return err
	}
	res.Close()

	if res.CommandTag().RowsAffected() == 0 {
		trans.Rollback(context.TODO())
		return ErrAddNotification
	}

	trans.Commit(context.TODO())
	return nil
}
