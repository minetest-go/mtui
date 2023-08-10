package web

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

/*
type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendError(w http.ResponseWriter, code int, message string) {
	logrus.WithFields(logrus.Fields{
		"code":    code,
		"message": message,
	}).Error("http error")

	errmsg := &ErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	buf := bytes.NewBuffer([]byte{})
	json.NewEncoder(buf).Encode(errmsg)
	w.Write(buf.Bytes())
}
*/

func SendError(w http.ResponseWriter, code int, message string) {
	logrus.WithFields(logrus.Fields{
		"code":    code,
		"message": message,
	}).Error("http error")
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func SendText(w http.ResponseWriter, txt string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(txt))
}

func SendJson(w http.ResponseWriter, o interface{}) []byte {
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

func Send(w http.ResponseWriter, o interface{}, err error) {
	if err != nil {
		SendError(w, 500, err.Error())
	} else {
		SendJson(w, o)
	}
}
