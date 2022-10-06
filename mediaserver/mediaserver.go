package mediaserver

type MediaServer struct {
	// file extensions
	extensions []string

	// sha1 to path
	Media map[string]string

	Size  int64
	Count int
}

func New() *MediaServer {
	return &MediaServer{
		extensions: []string{".png", ".jpg", ".jpeg", ".ogg", ".x", ".b3d", ".obj"},
		Media:      map[string]string{},
	}
}
