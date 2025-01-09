// handler/handler.go
package handle

import (
	"httpProject/session"
	"httpProject/stake"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	port          = 9000
	maxHighStakes = 20
)

type App struct {
	SessionManager *session.SessionManager
	StakeMap       *stake.StakeMap
	// dataStructure sync.Map // betOfferId -> datastructure
}

func NewApp() *App {
	return &App{
		SessionManager: session.NewSessionManager(),
		StakeMap:       stake.NewstakeMap(),
	}
}
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < 2 {
		app.sendResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	switch {
	case len(pathParts) == 2 && method == http.MethodGet && strings.HasSuffix(path, "/session"):
		app.handleGetSession(w, r, pathParts[0])
	case len(pathParts) == 2 && method == http.MethodGet && strings.HasSuffix(path, "/highstakes"):
		app.handleGetHighStakes(w, r, pathParts[0])
	case len(pathParts) == 2 && method == http.MethodPost && strings.HasSuffix(path, "/stake"):
		app.handlePostStake(w, r, pathParts[0])
	default:
		app.sendResponse(w, http.StatusNotFound, "Not Found")
	}
}

// 处理 GET /<customerid>/session
func (app *App) handleGetSession(w http.ResponseWriter, r *http.Request, customerID string) {
	ID, err := strconv.Atoi(customerID)
	if err != nil {
		app.sendResponse(w, http.StatusBadRequest, "need input number")
		return
	}
	log.Printf("handle get session")
	session := app.SessionManager.GetSession(ID)
	log.Printf(" get success")

	app.sendResponse(w, http.StatusOK, session.SessionKey)
}

type PostStakeRequest struct {
	Stake int `json:"stake"`
}

// 处理 POST /<betofferid>/stake?sessionkey=<sessionkey>
func (app *App) handlePostStake(w http.ResponseWriter, r *http.Request, betOfferIDstring string) {
	betOfferID, err := strconv.Atoi(betOfferIDstring)
	if err != nil {
		app.sendResponse(w, http.StatusBadRequest, "need input number")
		return
	}

	sessionKey := r.URL.Query().Get("session")
	if sessionKey == "" {
		app.sendResponse(w, http.StatusUnauthorized, "Session key required")
		return
	}

	customerID, ok := app.SessionManager.GetCustomerID(sessionKey)
	if !ok {
		app.sendResponse(w, http.StatusUnauthorized, "Invalid session key")
		return
	}

	// var stakeRequest PostStakeRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.sendResponse(w, http.StatusBadRequest, "Invalid input body")
		return
	}

	stakeStr := string(body)
	stake, err := strconv.Atoi(stakeStr) // 将字符串转成 int
	if err != nil {
		app.sendResponse(w, http.StatusBadRequest, "Invalid stake value")
		return
	}

	log.Printf("handle post stake")
	app.StakeMap.Insert(customerID, betOfferID, stake, maxHighStakes)

	app.sendResponse(w, http.StatusNoContent, "")
}

// 处理 GET /<betofferid>/highstakes
func (app *App) handleGetHighStakes(w http.ResponseWriter, r *http.Request, betOfferID string) {
	log.Printf(" start handle high stake")
	ID, err := strconv.Atoi(betOfferID) // 将字符串转成 int
	if err != nil {
		app.sendResponse(w, http.StatusBadRequest, "Invalid input betOfferID")
		return
	}
	topStakes, ok := app.StakeMap.GetTop(ID, 20)
	log.Printf(" handle highstake %d", len(topStakes))

	if !ok {
		app.sendResponse(w, http.StatusOK, "")
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strings.Join(topStakes, ",")))
}

func (app *App) sendResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain") // 设置 Content-Type 为 text/plain
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
