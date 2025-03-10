package mediaserver

import "go.uber.org/atomic"

var MEDIA_HEADER []byte = []byte("MTHS")
var MEDIA_VERSION = []byte{0, 1}

type MediaServer struct {
	// file extensions
	extensions []string

	// sha1 to path
	Media map[string]string

	Size             int64
	Count            int
	TransferredBytes atomic.Int64
}

func New() *MediaServer {
	return &MediaServer{
		extensions: []string{".png", ".jpg", ".jpeg", ".ogg", ".x", ".b3d", ".obj"},
		Media:      map[string]string{},
	}
}
