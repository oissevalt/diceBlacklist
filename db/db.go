package db

import (
	"database/sql"
	"diceBlacklist/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	initTable = `
	CREATE TABLE IF NOT EXISTS blacklist(
	    id TEXT PRIMARY KEY NOT NULL,
	    reason TEXT,
	    initial INTEGER NOT NULL,
	    latest INTEGER
	)
	`
	queryTable = `
	SELECT * FROM blacklist WHERE id = ?
	`
	insertRow = `
	INSERT INTO blacklist (id, reason, initial, latest) VALUES (?, ?, ?, ?)
	`
	updateRow = `
	UPDATE blacklist SET reason = ?, latest = ? WHERE id = ?
	`
	deleteRow = `
	DELETE FROM blacklist WHERE id = ?
	`
)

var (
	Database *sql.DB
	clientId []string
	Running  = new(atomic.Bool)
)

func init() {
	clientId = []string{}
	Running.Store(false)
}

type BlacklistItem struct {
	Id     string   `json:"id"`
	Reason []string `json:"reason"`
	First  int64    `json:"initial"`
	Last   int64    `json:"latest"`
}

func Initialize(interval int) error {
	err := ReadAppID()
	if err != nil {
		return err
	}

	Database, err = sql.Open("sqlite3", "blacklist.sqlite.db")
	if err != nil {
		return err
	}
	Running.Store(true)

	_, err = Database.Exec(initTable)
	if err != nil {
		return err
	}

	go backup(24 * time.Hour)
	return nil
}

func Authenticate(id string) bool {
	logger.Logger.Debugf("authenticating client %s", id)
	for _, i := range clientId {
		if i == id {
			return true
		}
	}
	return false
}

func Query(id string) (*BlacklistItem, error) {
	var item = new(BlacklistItem)
	var reasons string

	logger.Logger.Debugf("querying for %s", id)
	stmt, err := Database.Prepare(queryTable)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&item.Id, &reasons, &item.First, &item.Last)
	if err != nil {
		return nil, err
	}

	if reasons != "" {
		err = json.Unmarshal([]byte(reasons), &item.Reason)
		if err != nil {
			return nil, err
		}
	}

	return item, nil
}

func Add(id, reason string, timestamp int64) error {
	res, err := Query(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r := fmt.Sprintf(`["%s"]`, reason)
			stmt, err2 := Database.Prepare(insertRow)
			if err2 != nil {
				return err2
			}
			defer stmt.Close()

			_, err = stmt.Exec(id, r, timestamp, timestamp)
			return err
		} else {
			return err
		}
	}

	if len(res.Reason) <= 0 || (len(res.Reason) > 0 && res.Reason[len(res.Reason)-1] != reason) {
		res.Reason = append(res.Reason, reason)
	}
	res.Last = timestamp
	r, err := json.Marshal(res.Reason)
	if err != nil {
		return err
	}

	stmt, err := Database.Prepare(updateRow)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r, res.Last, id)
	if err != nil {
		return err
	}

	return nil
}

func Remove(id string) (error, bool) {
	_, err := Query(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false
		}
		return err, false
	}

	stmt, err := Database.Prepare(deleteRow)
	if err != nil {
		return err, false
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err, false
	}

	return nil, true
}

func backup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for tick := range ticker.C {
		logger.Logger.Info("start backup, database down")

		if s, err := os.Stat("backups"); err != nil || (s != nil && !s.IsDir()) {
			if err = os.MkdirAll("backups", 0755); err != nil {
				logger.Logger.Errorf("backup error: %v\n", err)
				continue
			}
		}

		func() {
			Database.Close()
			defer func() {
				logger.Logger.Debug("restarting database")

				var err error
				Database, err = sql.Open("sqlite3", "blacklist.sqlite.db")
				if err != nil {
					logger.Logger.Fatalf("database reboot error: %v", err)
				}
				Running.Store(true)

				logger.Logger.Info("database re-opened")
			}()

			Running.Store(false)
			t := tick.Format("2006-01-02_150405")
			dst := fmt.Sprintf("backups/blacklist_%s.sqlite.db", t)

			out, err := os.Open("blacklist.sqlite.db")
			if err != nil {
				logger.Logger.Errorf("backup error: %v", err)
				return
			}
			defer out.Close()

			in, err := os.Create(dst)
			if err != nil {
				logger.Logger.Errorf("backup error: %v", err)
				return
			}
			defer in.Sync()

			_, err = io.Copy(in, out)
			if err != nil {
				logger.Logger.Errorf("backup error: %v", err)
				return
			}
		}()
	}
}
