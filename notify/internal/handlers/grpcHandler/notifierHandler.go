package grpchandler

import (
	"context"
	"time"

	"github.com/ilovepitsa/happy/notify/api/notifier"
	"github.com/ilovepitsa/happy/notify/internal/entity"
)

type (
	NotificationStorage interface {
		GetUserNotifications(userId int32, idAbove int) ([]entity.Notification, int, error)
		GetUnsendEmailNotifications() ([]entity.Notification, error)
		GetNotificationText(isEmail bool) (string, error)
		AddNotification(target int32, isEmail bool, text string, date time.Time, notify_before time.Duration) error
	}

	NotificationServer struct {
		notifier.UnimplementedNotifierServer
		st NotificationStorage
	}
)

func NewNotificationServer(st NotificationStorage) *NotificationServer {
	return &NotificationServer{
		st: st,
	}
}

func (ns *NotificationServer) Create(ctx context.Context, info *notifier.NotificationInfo) (*notifier.Result, error) {
	date, _ := time.Parse("2006-01-02", info.Date)
	notify_before, _ := time.ParseDuration(info.NotifyBefore)
	textWeb, err := ns.getNotificationText(false)
	if err != nil {
		return &notifier.Result{Success: false}, err
	}
	textEmail, err := ns.getNotificationText(true)

	if err != nil {
		return &notifier.Result{Success: false}, err
	}

	err = ns.st.AddNotification(info.User.UserId, false, textWeb, date, notify_before)
	if err != nil {
		return &notifier.Result{Success: false}, err
	}
	err = ns.st.AddNotification(info.User.UserId, true, textEmail, date, notify_before)
	if err != nil {
		return &notifier.Result{Success: false}, err
	}

	return &notifier.Result{Success: true}, nil

}

func (ns *NotificationServer) getNotificationText(isEmail bool) (string, error) {
	return ns.st.GetNotificationText(isEmail)
}

func (ns *NotificationServer) GetUserNotifications(user *notifier.User, stream notifier.Notifier_GetUserNotificationsServer) error {

	lastId := 0
	var notifications []entity.Notification
	var err error
	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			notifications, lastId, err = ns.st.GetUserNotifications(user.UserId, lastId)
			if err != nil {
				return err
			}

			for _, notification := range notifications {
				err = stream.Send(&notifier.Notification{NotificationText: notification.NotificationText, Date: notification.Date.String()})
				if err != nil {
					return err
				}
			}
		}
	}
}
