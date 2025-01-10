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
	SessionsByCustomerID sync.Map // customerId -> *Session
	SessionsBySessionKey sync.Map // sessionKey -> string
	mu                   sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		SessionsByCustomerID: sync.Map{},
		SessionsBySessionKey: sync.Map{},
	}
}
func (m *SessionManager) GetSession(customerID int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	sessionValue, ok := m.SessionsByCustomerID.Load(customerID)
	log.Printf(" get session  %t in sessionmap,", ok)
	if ok {
		session := sessionValue.(*Session)
		if session.IsExpired() {
			log.Printf(" expire")

			m.SessionsByCustomerID.Delete(customerID)
			m.SessionsBySessionKey.Delete(session.SessionKey)
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
	m.SessionsByCustomerID.Store(customerID, session)
	m.SessionsBySessionKey.Store(session.SessionKey, customerID)

}

func (m *SessionManager) GetCustomerID(sessionKey string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	customerIDValue, ok := m.SessionsBySessionKey.Load(sessionKey)
	if !ok {
		return -1, false
	}
	sessionValue, ok := m.SessionsByCustomerID.Load(customerIDValue)
	if ok {
		session := sessionValue.(*Session)
		if session.IsExpired() {
			log.Printf(" this session key is expire")
			return -1, false
		}

	}

	return customerIDValue.(int), true
}

func (m *SessionManager) SessionCleanup() {
	log.Printf("Session cleanup started.")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		log.Printf("Waiting for next tick...")
		<-ticker.C
		log.Printf("session clean here")

		// Session cleanup logic
		m.SessionsByCustomerID.Range(func(key, value interface{}) bool {
			log.Printf("session clean here start range %d", key)

			session := value.(*Session)
			if session.IsExpired() {
				m.SessionsByCustomerID.Delete(key)
				m.SessionsBySessionKey.Delete(session.SessionKey)
				log.Printf("Deleted session for customer: %v\n", key)
			}
			return true
		})
	}
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
