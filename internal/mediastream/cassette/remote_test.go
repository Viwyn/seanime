package cassette

import (
	"seanime/internal/mediastream/videofile"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFixedKeyframeIndex(t *testing.T) {
	ki := newFixedKeyframeIndex("hash", 7)

	length, done := ki.Length()
	assert.True(t, done)
	assert.Equal(t, int32(4), length)
	assert.Equal(t, []float64{0, 2, 4, 6}, ki.Slice(0, length))
}

func TestBuildQualityLadderRemoteNeverTransmuxes(t *testing.T) {
	codec := "avc1.640028"
	video := &videofile.Video{MimeCodec: &codec, Width: 1920, Height: 1080}

	local := BuildQualityLadder(&videofile.MediaInfo{Path: "C:/anime/ep1.mkv", Video: video})
	assert.True(t, local[0].OriginalCanTransmux)

	remote := BuildQualityLadder(&videofile.MediaInfo{Path: "http://127.0.0.1:43211/api/v1/torrentstream/stream/ep1.mkv", Video: video})
	assert.False(t, remote[0].OriginalCanTransmux)
	assert.True(t, remote[0].NeedsTranscode)
}
