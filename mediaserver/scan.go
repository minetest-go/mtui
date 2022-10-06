package mediaserver

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// source: https://github.com/minetest-tools/mtmediasrv/blob/master/main.go#L131
func getHash(path string) (string, error) {
	var hashStr string

	f, err := os.Open(path)
	if err != nil {
		return hashStr, err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return hashStr, err
	}

	hashStr = hex.EncodeToString(h.Sum(nil)[:20])
	return hashStr, nil
}

func (m *MediaServer) ScanDefaultSubdirs(wd string) error {
	for _, dir := range []string{"game", "worldmods", "textures"} {
		fullpath := path.Join(wd, dir)
		info, err := os.Stat(fullpath)
		if err != nil || info == nil || !info.IsDir() {
			continue
		}

		err = m.Scan(fullpath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MediaServer) Scan(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, suffix := range m.extensions {
			if strings.HasSuffix(strings.ToLower(path), suffix) {
				hash, err := getHash(path)
				if err != nil {
					return err
				}
				if m.Media[hash] == "" {
					m.Count++
					m.Size += info.Size()
					m.Media[hash] = path
				}
			}
		}
		return nil
	})
	return err
}
