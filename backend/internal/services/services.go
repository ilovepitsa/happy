package services

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilovepitsa/happy/backend/internal/handler/auth"
)

type (
	NotifyInterface interface {
		GetSubscriber() UserInterface
		GetBirthdayBoy() UserInterface
		GetNotifyBefore() time.Duration
	}

	UserInterface interface {
		GetID() uint32
	}

	Auth interface {
		Check(c *gin.Context) (*auth.Session, error)
		Create(c *gin.Context, user UserInterface) error
		Delete(c *gin.Context) error
	}

	Notify interface {
		AddFollow(subscriber UserInterface, birthdayBoy UserInterface) error
	}

	Services struct {
		Auth   Auth
		Notify Notify
	}
)
