package dataType

import "time"

// Session 结构体
type Session struct {
	SessionKey string    `json:"session_key"`
	ExpiryTime time.Time `json:"expiry_time"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiryTime)
}
