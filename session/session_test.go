package session

import (
	"sync"
	"testing"
	"time"
)

func TestGetSession(t *testing.T) {
	sessionManager := NewSessionManager()

	// 测试获取不存在的 session
	session := sessionManager.GetSession(1)
	if session == nil {
		t.Errorf("Expected a new session, got nil")
	}

	// 测试获取已存在的 session
	session1 := sessionManager.GetSession(1)
	if session1 != session {
		t.Errorf("Expected the same session, got different sessions")
	}

	// 测试获取过期的 session
	session.ExpiryTime = time.Now().Add(-time.Minute)
	session2 := sessionManager.GetSession(1)
	if session2 == session {
		t.Errorf("Expected a new session, got the same expired session")
	}

	// 测试并发获取 session
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(3 * time.Second)
			sessionManager.GetSession(id)
		}(i)
	}
	wg.Wait()

}
