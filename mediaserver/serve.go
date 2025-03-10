package mediaserver

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func sendIndex(w http.ResponseWriter, hashes [][]byte) {
	headers := w.Header()
	headers.Add("Content-Type", "octet/stream")
	headers.Add("Content-Length", fmt.Sprintf("%d", 6+(len(hashes)*20)))

	w.Write(MEDIA_HEADER)
	w.Write(MEDIA_VERSION)
	for _, hash := range hashes {
		w.Write(hash)
	}
}

// adopted from: https://github.com/minetest-tools/mtmediasrv/blob/master/main.go#L58
func (m *MediaServer) ServeHTTPIndex(w http.ResponseWriter, r *http.Request) {

	// hashes to send
	hashes := [][]byte{}

	if r.Method == http.MethodGet {
		// simple GET send everything
		for hash := range m.Media {
			b, _ := hex.DecodeString(hash)
			hashes = append(hashes, b)
		}

		sendIndex(w, hashes)
		return
	}

	// only POST allowed from here on
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// check sent version
	header := make([]byte, 4)
	version := make([]byte, 2)
	r.Body.Read(header)
	r.Body.Read(version)

	if !bytes.Equal(header, MEDIA_HEADER) {
		http.Error(w, "invalid MTHS header", http.StatusInternalServerError)
		return
	}
	if !bytes.Equal(version, MEDIA_VERSION) {
		http.Error(w, "unsupported MTHS version", http.StatusInternalServerError)
		return
	}

	// read client needed hashes
	clientarr := make([]string, 0)
	for {
		h := make([]byte, 20)
		count, err := r.Body.Read(h)
		if err != nil && count == 0 {
			break
		}
		str := hex.EncodeToString(h)
		clientarr = append(clientarr, str)
	}
	r.Body.Close()

	for _, v := range clientarr {
		if m.Media[v] != "" {
			// we have the media for the hash
			b, _ := hex.DecodeString(v)
			hashes = append(hashes, b)
		}
	}

	sendIndex(w, hashes)
}

func (m *MediaServer) ServeHTTPFetch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	path := m.Media[hash]
	if path == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if file == nil {
		http.Error(w, "file vanished", http.StatusNotFound)
		return
	}

	info, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	headers := w.Header()
	headers.Add("Content-Type", "octet/stream")
	headers.Add("Content-Length", fmt.Sprintf("%d", info.Size()))

	// update stats
	m.TransferredBytes.Add(info.Size())

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
