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

// adopted from: https://github.com/minetest-tools/mtmediasrv/blob/master/main.go#L58
func (m *MediaServer) ServeHTTPIndex(w http.ResponseWriter, r *http.Request) {
	header := make([]byte, 4)
	version := make([]byte, 2)

	r.Body.Read(header)
	r.Body.Read(version)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !bytes.Equal(header, []byte("MTHS")) {
		http.Error(w, "invalid MTHS header", http.StatusInternalServerError)
		return
	}
	if !bytes.Equal(version, []byte{0, 1}) {
		http.Error(w, "unsupported MTHS version", http.StatusInternalServerError)
		return
	}

	// read client needed hashes
	clientarr := make([]string, 0)
	for {
		h := make([]byte, 20)
		_, err := r.Body.Read(h)
		if err != nil {
			break
		}
		clientarr = append(clientarr, hex.EncodeToString(h))
	}
	r.Body.Close()

	// Iterate over client hashes and remove hashes that we don't have from it
	resultmap := map[string]bool{}
	for _, v := range clientarr {
		if m.Media[v] != "" {
			resultmap[v] = true
		}
	}

	// formulate response
	headers := w.Header()
	headers.Add("Content-Type", "octet/stream")
	headers.Add("Content-Length", fmt.Sprintf("%d", 6+(len(resultmap)*20)))

	w.Write([]byte(header))
	w.Write([]byte(version))
	for k := range resultmap {
		b, _ := hex.DecodeString(k)
		w.Write([]byte(b))
	}
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
