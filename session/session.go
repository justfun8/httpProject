package session

import (
	"math/rand"
	"sync"
	"time"

	"log"
)

const sessionTimeoutMins = 10

type Session struct {
	SessionKey string    `json:"session_key"`
	ExpiryTime time.Time `json:"expiry_time"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiryTime)
}

type SessionManager struct {
	sessionsByCustomerID sync.Map // customerId -> *Session
	sessionsBySessionKey sync.Map // sessionKey -> string
	mu                   sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessionsByCustomerID: sync.Map{},
		sessionsBySessionKey: sync.Map{},
	}
}
func (m *SessionManager) GetSession(customerID int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	sessionValue, ok := m.sessionsByCustomerID.Load(customerID)
	log.Printf(" get session %t", ok)
	if ok {
		session := sessionValue.(*Session)
		if session.IsExpired() {
			log.Printf(" expire")

			m.sessionsByCustomerID.Delete(customerID)
			m.sessionsBySessionKey.Delete(session.SessionKey)
		} else {
			log.Printf(" not expire")

			return session
		}
	}
	newSession := &Session{
		SessionKey: m.generateSessionKey(),
		ExpiryTime: time.Now().Add(time.Duration(sessionTimeoutMins) * time.Second),
	}
	log.Printf(" key %s", newSession.SessionKey)

	m.storeSession(customerID, newSession)
	log.Printf(" susccess")

	return newSession
}

func (m *SessionManager) storeSession(customerID int, session *Session) {
	m.sessionsByCustomerID.Store(customerID, session)
	m.sessionsBySessionKey.Store(session.SessionKey, customerID)
}

func (m *SessionManager) GetCustomerID(sessionKey string) (int, bool) {
	customerIDValue, ok := m.sessionsBySessionKey.Load(sessionKey)
	if !ok {
		return -1, false
	}
	return customerIDValue.(int), true
}

func (m *SessionManager) SessionCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	m.sessionsByCustomerID.Range(func(key, value interface{}) bool {
		session := value.(*Session)
		if session.IsExpired() {
			m.sessionsByCustomerID.Delete(key)
			m.sessionsBySessionKey.Delete(session.SessionKey)
			log.Printf("Deleted session for customer: %v\n", key)
		}
		return true
	})
}

func (m *SessionManager) generateSessionKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 20)
	for i := range b {
		const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}
