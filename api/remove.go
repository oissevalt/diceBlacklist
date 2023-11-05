package api

import (
	"diceBlacklist/db"
	"diceBlacklist/logger"
	"net/http"
)

func Remove(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info("received incoming removal request")

	client := r.Header.Get("appid")
	if client == "" || !db.Authenticate(client) {
		logger.Logger.Info("client is not permitted to remove items")
		Respond(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, _ := Sanitize(r.URL.Query().Get("id"))
	if id == "" || r.Method != http.MethodDelete {
		Respond(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err, found := db.Remove(id)
	if err != nil { // Remove already sorted out ErrNoRows
		logger.Logger.Errorf("failed to remove item %s: %v", id, err)
		Respond(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !found {
		logger.Logger.Infof("could not find %s to remove", id)
		Respond(w, "Not Found", http.StatusNotFound)
		return
	}

	logger.Logger.Infof("removed item with ID %s", id)
}
