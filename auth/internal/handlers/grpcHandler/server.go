package grpchandler

import (
	"context"

	"github.com/ilovepitsa/happy/auth/api/sessions"
)

type SessionManager interface {
	Check(*sessions.SessionID) (*sessions.Session, error)
	Create(*sessions.Session) (*sessions.SessionID, error)
	Delete(*sessions.SessionID) error
}

type SessionHandler struct {
	sessions.UnimplementedAuthCheckerServer
	manager SessionManager
}

func NewSessionHandler(manager SessionManager) *SessionHandler {
	return &SessionHandler{
		manager: manager,
	}
}
func (s *SessionHandler) Create(ctx context.Context, sess *sessions.Session) (*sessions.SessionID, error) {
	return s.manager.Create(sess)
}
func (s *SessionHandler) Check(ctx context.Context, sessId *sessions.SessionID) (*sessions.Session, error) {
	return s.manager.Check(sessId)
}
func (s *SessionHandler) Delete(ctx context.Context, sessId *sessions.SessionID) (*sessions.Nothing, error) {
	return &sessions.Nothing{}, s.manager.Delete(sessId)
}
