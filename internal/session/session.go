package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

type Session struct {
	UID         int
	SessionID   string
	CreatedAt   time.Time
	LastUpdated time.Time
}

type SessionManager struct {
	sessionMap      map[string]*Session
	userMap         map[int]string
	mutex           sync.RWMutex
	sessionDuration time.Duration
	cleanupInterval time.Duration
}

func NewSessionManager(sessionDuration, cleanupInterval time.Duration) *SessionManager {
	manager := &SessionManager{
		sessionMap:      make(map[string]*Session),
		userMap:         make(map[int]string),
		sessionDuration: sessionDuration,
		cleanupInterval: cleanupInterval,
	}
	go manager.periodicCleanup()
	return manager
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (sm *SessionManager) GetUIDBySessionID(session_id string) int {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	s, ok := sm.sessionMap[session_id]
	if !ok {
		return -1
	}

	return s.UID
}

func (sm *SessionManager) CreateOrUpdateSession(uid int) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if old_sess_id, exists := sm.userMap[uid]; exists {
		delete(sm.sessionMap, old_sess_id)
	}

	sess_id, err := generateSessionID()
	if err != nil {
		return "", err
	}

	now := time.Now()
	session := &Session{
		UID:         uid,
		SessionID:   sess_id,
		CreatedAt:   now,
		LastUpdated: now,
	}

	sm.sessionMap[sess_id] = session
	sm.userMap[uid] = sess_id
	fmt.Printf("Session created: uid=%d, session_id=%s\n", uid, sess_id)

	return sess_id, nil
}

func (sm *SessionManager) GetSessionBySessionID(session_id string) (*Session, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	session, exists := sm.sessionMap[session_id]
	return session, exists
}

func (sm *SessionManager) GetSessionIDByUID(uid int) (string, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	sessionID, exists := sm.userMap[uid]
	return sessionID, exists
}

func (sm *SessionManager) DeleteSessionByUID(uid int) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if session_id, exists := sm.userMap[uid]; exists {
		delete(sm.sessionMap, session_id)
		delete(sm.userMap, uid)
		fmt.Printf("Session deleted: uid=%d, session_id=%s\n", uid, session_id)
	}
}

func (sm *SessionManager) DeleteSessionBySessionID(session_id string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if session, exists := sm.sessionMap[session_id]; exists {
		delete(sm.userMap, session.UID)
		delete(sm.sessionMap, session_id)
		fmt.Printf("Session deleted: uid=%d, session_id=%s\n", session.UID, session_id)
	}
}

func (sm *SessionManager) CheckSessionExists(session_id string) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.sessionMap[session_id]; exists {
		return true
	}

	return false
}

func (sm *SessionManager) periodicCleanup() {
	ticker := time.NewTicker(sm.cleanupInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		sm.cleanupExpiredSessions()
	}
}

func (sm *SessionManager) cleanupExpiredSessions() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	now := time.Now()
	for sessionID, session := range sm.sessionMap {
		if now.Sub(session.LastUpdated) > sm.sessionDuration {
			delete(sm.sessionMap, sessionID)
			delete(sm.userMap, session.UID)
			fmt.Printf("Session is cleaned: uid=%d, sessionID=%s\n", session.UID, sessionID)
		}
	}
}
