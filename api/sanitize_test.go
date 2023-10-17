package api

import "testing"

func TestSanitize(t *testing.T) {
	var ids = []string{
		"QQ:12345",
		"QQ12345",
		"QQ:123ab",
		"QQ-CH-Group:12345",
		"QQ-Group:12345",
		"QQ-Ch-groUp:12345",
		"QQ-CH:12345",
		"disCord:12345",
		"Kook: 12345",
		"Dis-chan-gp:12345",
	}

	for _, id := range ids {
		sanitized, ok := Sanitize(id)
		t.Logf("Original: %18s Sanitized: %18s OK: %v", id, sanitized, ok)
	}
}
