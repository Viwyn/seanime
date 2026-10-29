package videofile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHashForStream(t *testing.T) {
	a := GetHashForStream("abc/ep1.mkv")
	assert.Len(t, a, 40)
	assert.Equal(t, a, GetHashForStream("abc/ep1.mkv"))
	assert.NotEqual(t, a, GetHashForStream("abc/ep2.mkv"))
}

func TestExtensionFromPath(t *testing.T) {
	assert.Equal(t, "mkv", extensionFromPath("C:/anime/ep1.mkv"))
	assert.Equal(t, "mkv", extensionFromPath("http://127.0.0.1:43211/api/v1/torrentstream/stream/%5BGroup%5D%20Ep%201.mkv?token=a.b.c"))
	assert.Equal(t, "", extensionFromPath("http://127.0.0.1:43211/api/v1/torrentstream/stream/noext"))
}
