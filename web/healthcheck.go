package web

import (
	"net/http"
	"time"
)

var last_healthcheck_result error
var last_healthcheck_time int64

func (a *Api) HealthCheck(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Unix()
	if now-last_healthcheck_time > 10 {
		// last healthcheck was over 10 seconds ago, refresh
		last_healthcheck_result = a.app.Healthcheck()
		last_healthcheck_time = now
	}

	Send(w, true, last_healthcheck_result)
}
