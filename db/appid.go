package db

import (
	"diceBlacklist/logger"
	"encoding/json"
	"errors"
	"github.com/fsnotify/fsnotify"
	"os"
)

func InitAppID() (err error) {
	clientId = []string{}
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
			// Some editors use temp file, then it will trigger CHMOD, REMOVE and CREATE.
			// This makes things too complicated so it's ignored.
			if event.Has(fsnotify.Write) {
				if localErr := InitAppID(); localErr != nil {
					logger.Logger.Errorf("failed to reload appid list: %v", localErr)
					return
				}
				logger.Logger.Info("appid list reloaded")
			}
		case watcherErr, ok := <-w.Errors:
			if !ok {
				return
			}
			logger.Logger.Errorf("file watcher error: %v", watcherErr)
		}
	}
}