package main

import (
	"diceBlacklist/api"
	"diceBlacklist/db"
	"diceBlacklist/logger"
	"flag"
	"fmt"
	"net/http"
)

var port = 3518

func main() {
	logger.Logger.Debug("program started")
	flag.IntVar(&port, "port", 3518, "Server port number.")

	if err := db.InitDatabase(); err != nil {
		logger.Logger.Fatalf("database initialization error: %v", err)
	}

	if err := db.Watch(); err != nil {
		logger.Logger.Errorf("failed to initialize file watcher: %v", err)
	}

	http.HandleFunc("/query", api.Limit(api.Query))
	http.HandleFunc("/report", api.Limit(api.Report))
	http.HandleFunc("/remove", api.Limit(api.Remove))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})
	logger.Logger.Infof("server running on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Logger.Errorf("server error: %v", err)
	}
}
