package web

import (
	"encoding/json"
	"fmt"
	"io"
	"mtui/app"
	"mtui/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/sftp"
	"github.com/studio-b12/gowebdav"
	"golang.org/x/crypto/ssh"
)

type RestoreJobState string

const (
	RestoreJobRunning RestoreJobState = "running"
	RestoreJobSuccess RestoreJobState = "success"
	RestoreJobFailure RestoreJobState = "failure"
)

type RestoreJobInfo struct {
	ID      string          `json:"id"`
	Status  RestoreJobState `json:"state"`
	Message string          `json:"message"`
}

type RestoreJobType string

const (
	RestoreJobTypeSCP    RestoreJobType = "scp"
	RestoreJobTypeWEBDAV RestoreJobType = "webdav"
)

type CreateRestoreJob struct {
	ID       string         `json:"id"`
	Type     RestoreJobType `json:"type"`
	Host     string         `json:"host"` //scp
	Port     int            `json:"port"` //scp
	URL      string         `json:"url"`  // webdav
	Filename string         `json:"filename"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Key      string         `json:"key"`
}

var Restorejobs = map[string]*RestoreJobInfo{}

func restoreJob(a *app.App, job *CreateRestoreJob, info *RestoreJobInfo, c *types.Claims) {
	var reader io.Reader
	var err error

	switch job.Type {
	case RestoreJobTypeSCP:
		addr := fmt.Sprintf("%s:%d", job.Host, job.Port)
		config := &ssh.ClientConfig{
			User: job.Username,
			Auth: []ssh.AuthMethod{
				ssh.Password(job.Password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		client, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			info.Status = RestoreJobFailure
			info.Message = fmt.Sprintf("ssh dial failed: %v", err)
			return
		}
		defer client.Close()

		sc, err := sftp.NewClient(client)
		if err != nil {
			info.Status = RestoreJobFailure
			info.Message = fmt.Sprintf("sftp open failed: %v", err)
			return
		}
		defer sc.Close()

		file, err := sc.Open(job.Filename)
		if err != nil {
			info.Status = RestoreJobFailure
			info.Message = fmt.Sprintf("sftp create failed: %v", err)
			return
		}
		defer file.Close()
		reader = file

	case RestoreJobTypeWEBDAV:
		c := gowebdav.NewClient(job.URL, job.Username, job.Password)
		err = c.Connect()
		if err != nil {
			info.Status = RestoreJobFailure
			info.Message = fmt.Sprintf("webdav connection failed: %v", err)
			return
		}

		fi, err := c.Stat(job.Filename)
		if err != nil {
			info.Status = RestoreJobFailure
			info.Message = fmt.Sprintf("stat connection failed: %v", err)
			return
		}
		info.Message = fmt.Sprintf("Downlading zip file, size: %d bytes", fi.Size())

		r, err := c.ReadStream(job.Filename)
		if err != nil {
			info.Status = RestoreJobFailure
			info.Message = fmt.Sprintf("webdav connection failed: %v", err)
			return
		}
		defer r.Close()
		reader = r

	default:
		info.Status = RestoreJobFailure
		info.Message = fmt.Sprintf("unknown job type: %s", job.Type)
		return
	}

	if job.Key != "" {
		// enable decryption
		reader, err = app.EncryptedReader(job.Key, reader)
		if err != nil {
			info.Status = RestoreJobFailure
			info.Message = fmt.Sprintf("decryption failed: %v", err)
			return
		}
	}

	filecount := 0
	bytes, err := a.DownloadAndUnzip(a.WorldDir, reader, nil, c, &app.DownloadZipOpts{
		Callback: func(files, bytes int64, currentfile string) {
			info.Message = fmt.Sprintf("Copying file '%s' (progress: %d bytes, %d files)", currentfile, bytes, files)
			filecount++
		},
	})
	if err != nil {
		info.Status = RestoreJobFailure
		info.Message = fmt.Sprintf("restore failed: %v", err)
		return
	}

	// reinstall system-mods
	err = a.ReconfigureSystemMods()
	if err != nil {
		info.Status = RestoreJobFailure
		info.Message = fmt.Sprintf("mod reconfiguration failed: %v", err)
		return
	}

	info.Message = fmt.Sprintf("Restore complete with %d bytes and %d files", bytes, filecount)
	info.Status = RestoreJobSuccess

	a.CreateUILogEntry(&types.Log{
		Username: c.Username,
		Event:    "backup",
		Message:  info.Message,
	}, nil)
}

func (a *Api) CreateRestoreJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	job := &CreateRestoreJob{}
	err := json.NewDecoder(r.Body).Decode(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json decode error: %v", err))
		return
	}

	if job.ID == "" {
		job.ID = uuid.NewString()
	}

	info := &RestoreJobInfo{
		Status: RestoreJobRunning,
		ID:     job.ID,
	}
	Restorejobs[job.ID] = info
	go restoreJob(a.app, job, info, c)

	SendJson(w, info)
}

func (a *Api) GetRestoreJobInfo(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	id := mux.Vars(r)["id"]
	info := Restorejobs[id]
	if info == nil {
		SendError(w, 404, fmt.Errorf("restore job not found: %s", id))
		return
	}
	SendJson(w, info)
}
