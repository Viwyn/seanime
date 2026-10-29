package videofile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// GetHashFromPath returns a deterministic hash derived from the file path and
// its last-modified timestamp. The hash changes when the file is re-written,
// ensuring stale cache entries are invalidated.
//
// Uses SHA-256 with a stable time format (UnixNano) to avoid locale-dependent
// string representations of ModTime.
func GetHashFromPath(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	// Use UnixNano instead of ModTime().String() which can vary across
	// locales and Go versions.
	_, _ = fmt.Fprintf(h, "%s:%d", path, info.ModTime().UnixNano())
	return hex.EncodeToString(h.Sum(nil))[:40], nil // 40-char hex prefix
}

// GetHashForStream returns a deterministic hash for a remote stream (e.g. a torrent file)
// that has no stable path or modification time on disk.
func GetHashForStream(id string) string {
	h := sha256.Sum256([]byte("stream:" + id))
	return hex.EncodeToString(h[:])[:40]
}

// IsRemotePath reports whether the path is an HTTP(S) URL rather than a file on disk.
func IsRemotePath(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// extensionFromPath returns the extension without the dot, ignoring any URL query string.
func extensionFromPath(path string) string {
	if IsRemotePath(path) {
		if u, err := url.Parse(path); err == nil {
			path = u.Path
		}
	}
	return strings.TrimPrefix(filepath.Ext(path), ".")
}
