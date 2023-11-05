package api

import (
	"diceBlacklist/logger"
	"net/http"
	"strconv"
	"strings"
)

func Respond(w http.ResponseWriter, message string, status int) {
	http.Error(w, message, status)
	logger.Logger.Infof("server responded with HTTP %d: %s", status, message)
}

// Sanitize formats the provided id by trimming and adjusting aliases.
// It returns the sanitized string and true, or, if the id is empty or
// badly formed, the original string and false.
func Sanitize(id string) (string, bool) {
	if id == "" {
		return id, false
	}

	trimmed := strings.TrimSpace(id)
	seg := strings.Split(trimmed, ":")
	if len(seg) != 2 {
		return id, false
	}

	num := strings.TrimSpace(seg[1])
	if _, err := strconv.Atoi(num); err != nil {
		return id, false
	}

	idt := strings.SplitN(seg[0], "-", 3)
	plat := canonicalize(idt[0])
	if len(idt) > 1 {
		if len(idt) == 2 && strings.EqualFold(idt[1], "Group") {
			return plat + "-Group:" + num, true
		} else if len(idt) == 3 && strings.EqualFold(idt[1], "CH") && strings.EqualFold(idt[2], "Group") {
			return plat + "-CH-Group:" + num, true
		}
		return id, false
	} else {
		return plat + ":" + num, true
	}
}

func canonicalize(pl string) string {
	pl = strings.ToUpper(pl)
	switch pl {
	case "DISCORD":
		return "DC"
	case "TELEGRAM":
		return "TG"
	case "MINECRAFT":
		return "MC"
	default:
		return pl
	}
}
