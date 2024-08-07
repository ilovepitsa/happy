package repo

import (
	"context"

	"github.com/ilovepitsa/happy/notify/internal/entity"
	"github.com/jackc/pgx/v5"
)

type NotificationRepo struct {
	connection *pgx.Conn
}

func NewNotificationRepo(connection *pgx.Conn) *NotificationRepo {
	return &NotificationRepo{
		connection: connection,
	}
}

func (nr *NotificationRepo) GetUserNotifications(userId int32) ([]entity.Notification, error) {
	trans, err := nr.connection.Begin(context.TODO())
	if err != nil {
		trans.Rollback(context.TODO())
		return nil, err
	}

	res, err := trans.Query(context.TODO(), "select target, text, date from notification where type = 'web' and target = '$1'", userId)
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

	return result, nil
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
		return "", err
	}
	return text, nil
}
