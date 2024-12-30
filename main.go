package main

import (
	"httpProject/dataType"
	// "main/dataType"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	port               = 9000
	sessionTimeoutMins = 10
	maxHighStakes      = 20
	sessionKeyLength   = 20
)

type App struct {
	sessions sync.Map // customerId -> session
	betMap   sync.Map // betOfferId -> linkList
	mu       sync.RWMutex
}

func NewApp() *App {
	return &App{
		sessions: sync.Map{},
		betMap:   sync.Map{},
	}
}
func main() {
	app := NewApp()

	// 启动会话清理任务
	go app.sessionCleanup()

	http.HandleFunc("/", app.handler)

	log.Printf("Server starting on port: %d\n", port)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}

}

func (app *App) handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < 2 {
		app.sendResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	switch {
	case len(pathParts) == 2 && method == http.MethodGet && strings.HasSuffix(path, "/session"):
		customerID, err := strconv.Atoi(pathParts[0])
		if err != nil {
			app.sendResponse(w, http.StatusBadRequest, "Invalid Customer ID")
			return
		}
		app.handleGetSession(w, r, customerID)
	case len(pathParts) == 2 && method == http.MethodGet && strings.HasSuffix(path, "/highstakes"):
		betID, err := strconv.Atoi(pathParts[0])
		if err != nil {
			app.sendResponse(w, http.StatusBadRequest, "Invalid bet ID")
			return
		}
		app.handleGetHighStakes(w, r, betID)
	case len(pathParts) == 2 && method == http.MethodPost && strings.HasSuffix(path, "/stake"):
		betID, err := strconv.Atoi(pathParts[0])
		if err != nil {
			app.sendResponse(w, http.StatusBadRequest, "Invalid bet ID")
			return
		}
		app.handlePostStake(w, r, betID)
	default:
		app.sendResponse(w, http.StatusNotFound, "Not Found")
	}
}

// 处理 GET /<customerid>/session
func (app *App) handleGetSession(w http.ResponseWriter, r *http.Request, customerID int) {
	session, ok := app.sessions.Load(customerID)
	if ok {
		s := session.(*dataType.Session)
		if !s.IsExpired() {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(s.SessionKey))
			// app.sendJSONResponse(w, http.StatusOK, s.SessionKey)
			return
		}
	}

	newSession := &dataType.Session{
		SessionKey: app.generateSessionKey(),
		ExpiryTime: time.Now().Add(time.Duration(sessionTimeoutMins) * time.Minute),
	}
	app.sessions.Store(customerID, newSession)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newSession.SessionKey))
}

type PostStakeRequest struct {
	Stake int `json:"stake"`
}

// 处理 POST /<betofferid>/stake?sessionkey=<sessionkey>
func (app *App) handlePostStake(w http.ResponseWriter, r *http.Request, betOfferID int) {
	sessionKey := r.URL.Query().Get("session")
	if sessionKey == "" {
		app.sendResponse(w, http.StatusUnauthorized, "Session key required")
		return
	}

	customerID, ok := app.getCustomerIdFromSessionKey(sessionKey)
	if !ok {
		app.sendResponse(w, http.StatusUnauthorized, "Invalid session key")
		return
	}

	var stakeRequest PostStakeRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.sendResponse(w, http.StatusBadRequest, "Invalid input body")
		return
	}

	if err := json.Unmarshal(body, &stakeRequest); err != nil {
		app.sendResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	stake := stakeRequest.Stake
	app.mu.Lock()
	defer app.mu.Unlock()
	oldlist, ok := app.betMap.Load(betOfferID) // 这里使用Load获取，如果不存在则创建

	if !ok {
		list := dataType.NewDoublyLinkedList(maxHighStakes)
		list.Insert(customerID, stake)
		log.Printf("addmain%d ", customerID)
		app.betMap.Store(betOfferID, list) // 使用 Store 添加数据
		// betMap = list
	} else {
		olist := oldlist.(*dataType.DoublyLinkedList) // 断言类型
		olist.Insert(customerID, stake)
		app.betMap.Store(betOfferID, olist) // 使用 Store 添加数据

	}

}

// 处理 GET /<betofferid>/highstakes
func (app *App) handleGetHighStakes(w http.ResponseWriter, r *http.Request, betOfferID int) {
	betMapValue, ok := app.betMap.Load(betOfferID)
	if !ok {
		app.sendResponse(w, http.StatusOK, "")
		return
	}
	linkList := betMapValue.(*dataType.DoublyLinkedList)
	topStakes := linkList.GetTop(maxHighStakes)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strings.Join(topStakes, ",")))
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (app *App) sendResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain") // 设置 Content-Type 为 text/plain
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func (app *App) generateSessionKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, sessionKeyLength)
	for i := range b {
		const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}

func (app *App) sessionCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		app.sessions.Range(func(key, value interface{}) bool {
			session := value.(*dataType.Session)
			if session.IsExpired() {
				app.sessions.Delete(key)
				log.Printf("Deleted session for customer: %v\n", key)
			}
			return true
		})
	}
}

func (app *App) getCustomerIdFromSessionKey(sessionKey string) (int, bool) {
	var customerId int
	find := false
	app.sessions.Range(func(key, value interface{}) bool {
		session := value.(*dataType.Session)
		log.Printf("Checking session key: %s, stored session key: %s, customerID: %d\n", sessionKey, session.SessionKey, key) // 添加日志
		if session.SessionKey == sessionKey {
			customerId = key.(int)
			log.Printf("Found customer id: %d for session key: %v\n", customerId, sessionKey)
			find = true
			return false // 找到了匹配项，停止遍历
		}
		return true // 继续遍历
	})
	if !find {
		log.Printf("Could not find a customer id for session key: %v\n", sessionKey)
	}
	return customerId, find
}
