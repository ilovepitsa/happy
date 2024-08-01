package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ilovepitsa/happy/backend/internal/handler/auth"
)

type (
	UserInterface interface {
		GetID() uint32
	}

	SessionManager interface {
		Check(c *gin.Context) (*auth.Session, error)
		Create(c *gin.Context, user UserInterface) error
		Delete(c *gin.Context) error
	}

	AuthMiddleware struct {
		sm SessionManager
	}
)

func NewAuthMiddleware(sm SessionManager) *AuthMiddleware {
	return &AuthMiddleware{
		sm: sm,
	}
}

func (a *AuthMiddleware) UserIndetify(c *gin.Context) {

	sess, err := a.sm.Check(c)

	switch err {
	case auth.ErrNoAuth:
		c.Redirect(http.StatusFound, "/singin")
	case auth.ErrAuthService:
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	default:
	}

	c.Set(strconv.Itoa(int(auth.SessionKey)), sess)
	c.Next()
}
