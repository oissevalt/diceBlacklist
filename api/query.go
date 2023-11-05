package api

import (
	"database/sql"
	"diceBlacklist/db"
	"diceBlacklist/logger"
	"encoding/json"
	"errors"
	"net/http"
)

func Query(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info("received incoming query request")

	id, _ := Sanitize(r.URL.Query().Get("id"))
	if id == "" || r.Method != http.MethodGet {
		Respond(w, "Bad Request", http.StatusBadRequest)
		return
	}

	res, err := db.Query(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Logger.Debugf("query for %s selected no rows", id)
			Respond(w, "Not Found", http.StatusNotFound)
			return
		}
		logger.Logger.Errorf("query for %s terminated with error: %v", id, err)
		Respond(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Logger.Debugf("queried data: %v", res)

	con, err := json.Marshal(*res)
	if err != nil {
		logger.Logger.Errorf("failed to marshal item %s: %v", id, err)
		Respond(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(con)
	if err != nil {
		logger.Logger.Errorf("failed to write response: %v", err)
		Respond(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	logger.Logger.Infof("server responded with query result of %s", id)
}
