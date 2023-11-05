package api

import (
	"diceBlacklist/db"
	"diceBlacklist/logger"
	"net/http"
	"time"
)

func Report(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info("received incoming report")

	client := r.Header.Get("appid")
	if client == "" || !db.Authenticate(client) {
		logger.Logger.Info("client is not permitted to add to blacklist")
		Respond(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, ok := Sanitize(r.URL.Query().Get("id"))
	reason := r.URL.Query().Get("for")
	if !ok || id == "" || reason == "" || r.Method != http.MethodPost {
		Respond(w, "Bad Request", http.StatusBadRequest)
		return
	}

	ts := time.Now().Unix()
	err := db.Add(id, reason, ts)
	if err != nil { // ErrNoRows excluded
		logger.Logger.Errorf("failed to add item %s: %v", id, err)
		Respond(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Logger.Infof("added item with ID %s: %v", id, reason)
}
