package web

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func SendError(w http.ResponseWriter, code int, err error) {
	logrus.WithFields(logrus.Fields{
		"code": code,
		"err":  err,
	}).Error("http error")
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

func SendText(w http.ResponseWriter, txt string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(txt))
}

func SendJson(w http.ResponseWriter, o any) []byte {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(o)
	w.Write(buf.Bytes())
	return buf.Bytes()
}

func SendRawJson(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func Send(w http.ResponseWriter, o any, err error) {
	if err != nil {
		SendError(w, 500, err)
	} else {
		SendJson(w, o)
	}
}
