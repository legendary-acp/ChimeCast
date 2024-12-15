package session

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserName  string
	UserID    string // Added UserID field
	ExpiresAt time.Time
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

// CreateSession now takes both userName and userID
func (sm *SessionManager) CreateSession(userName string, userID string) (string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionID := uuid.NewString()
	sm.sessions[sessionID] = &Session{
		UserName:  userName,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	return sessionID, nil
}

// GetSession remains the same
func (sm *SessionManager) GetSession(sessionID string) (*Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil, errors.New("invalid or expired session")
	}
	return session, nil
}

// DeleteSession remains the same
func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, sessionID)
}
