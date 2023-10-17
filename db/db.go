package db

import (
	"database/sql"
	"diceBlacklist/logger"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"slices"
)

const (
	initTable = `
	CREATE TABLE IF NOT EXISTS blacklist(
	    id TEXT PRIMARY KEY NOT NULL,
	    reason TEXT,
	    first INTEGER NOT NULL,
	    last INTEGER
	)
	`
	queryTable = `
	SELECT * FROM blacklist WHERE id = ?
	`
	insertRow = `
	INSERT INTO blacklist (id, reason, first, last) VALUES (?, ?, ?, ?)
	`
	updateRow = `
	UPDATE blacklist SET reason = ?, last = ? WHERE id = ?
	`
	deleteRow = `
	DELETE FROM blacklist WHERE id = ?
	`
)

var (
	Database *sql.DB
	clientId []string
)

type BlacklistItem struct {
	Id     string   `json:"id"`
	Reason []string `json:"reason"`
	First  int64    `json:"first"`
	Last   int64    `json:"last"`
}

func InitDatabase() error {
	err := InitAppID()
	if err != nil {
		return err
	}

	Database, err = sql.Open("sqlite3", "blacklist.sqlite.db")
	if err != nil {
		return err
	}
	_, err = Database.Exec(initTable)
	if err != nil {
		return err
	}

	return nil
}

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

func Authenticate(id string) bool {
	logger.Logger.Debugf("authenticating client %s", id)
	_, found := slices.BinarySearch(clientId, id)
	return found
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
			r := fmt.Sprintf("[%s]", reason)
			stmt, err2 := Database.Prepare(insertRow)
			if err2 != nil {
				return err2
			}
			defer stmt.Close()

			_, err = stmt.Exec(id, r, timestamp, timestamp)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	res.Reason = append(res.Reason, reason)
	res.Last = timestamp
	r, err := json.Marshal(res.Reason)

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
