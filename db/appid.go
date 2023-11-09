package db

import (
	"diceBlacklist/logger"
	"encoding/json"
	"errors"
	"github.com/fsnotify/fsnotify"
	"os"
)

func init() {
	clientId = []string{}
}

func ReadAppID() (err error) {
	con, err := os.ReadFile("appid.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Logger.Warn("`appid.json` not found, creating a new one")
			file, err2 := os.Create("appid.json")
			if err2 != nil {
				return err2
			}
			_, _ = file.WriteString("[]")
		} else {
			return err
		}
	}

	if con != nil && len(con) > 0 {
		err = json.Unmarshal(con, &clientId)
		if err != nil {
			return err
		}
	}

	return nil
}

func Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go watcherFunc(watcher)

	err = watcher.Add("appid.json")
	if err != nil {
		return err
	}

	return nil
}

func watcherFunc(w *fsnotify.Watcher) {
	logger.Logger.Info("file watcher up and running")
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			logger.Logger.Debugf(
				"file watcher detected changes to %s: %v",
				event.Name,
				event.Op)
			// Some editors create a temporary file, so they will trigger
			// CHMOD, REMOVE and CREATE at once. For the sake of simplicity
			// this is not planned to be handled.
			if event.Has(fsnotify.Write) {
				if err := ReadAppID(); err != nil {
					logger.Logger.Errorf("failed to reload appid list: %v", err)
					return
				}
				logger.Logger.Info("appid list reloaded")
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			logger.Logger.Errorf("file watcher error: %v", err)
		}
	}
}
