package web

import (
	"encoding/json"
	"fmt"
	"io"
	"mtui/app"
	"mtui/types"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/sftp"
	"github.com/studio-b12/gowebdav"
	"golang.org/x/crypto/ssh"
)

type BackupJobState string

const (
	BackupJobRunning BackupJobState = "running"
	BackupJobSuccess BackupJobState = "success"
	BackupJobFailure BackupJobState = "failure"
)

type BackupJobInfo struct {
	ID      string         `json:"id"`
	Status  BackupJobState `json:"state"`
	Message string         `json:"message"`
}

type BackupJobType string

const (
	BackupJobTypeSCP    BackupJobType = "scp"
	BackupJobTypeWEBDAV BackupJobType = "webdav"
)

type CreateBackupJob struct {
	ID       string        `json:"id"`
	Type     BackupJobType `json:"type"`
	Host     string        `json:"host"` //scp
	Port     int           `json:"port"` //scp
	URL      string        `json:"url"`  // webdav
	Filename string        `json:"filename"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Key      string        `json:"key"`
}

var backupjobs = map[string]*BackupJobInfo{}

func backupJob(a *app.App, job *CreateBackupJob, info *BackupJobInfo, c *types.Claims) {
	var output io.WriteCloser
	var err error

	switch job.Type {
	case BackupJobTypeSCP:
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
			info.Status = BackupJobFailure
			info.Message = fmt.Sprintf("ssh dial failed: %v", err)
			return
		}
		defer client.Close()

		sc, err := sftp.NewClient(client, sftp.UseConcurrentWrites(true))
		if err != nil {
			info.Status = BackupJobFailure
			info.Message = fmt.Sprintf("sftp open failed: %v", err)
			return
		}
		defer sc.Close()

		file, err := sc.OpenFile(job.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		if err != nil {
			info.Status = BackupJobFailure
			info.Message = fmt.Sprintf("sftp create failed: %v", err)
			return
		}
		defer file.Close()
		output = file

	case BackupJobTypeWEBDAV:
		c := gowebdav.NewClient(job.URL, job.Username, job.Password)
		err = c.Connect()
		if err != nil {
			info.Status = BackupJobFailure
			info.Message = fmt.Sprintf("webdav connection failed: %v", err)
			return
		}

		pr, pw := io.Pipe()
		output = pw

		go func() {
			err := c.WriteStream(job.Filename, pr, 0644)
			if err != nil {
				info.Status = BackupJobFailure
				info.Message = fmt.Sprintf("webdav stream failed: %v", err)
			}
		}()

	default:
		info.Status = BackupJobFailure
		info.Message = fmt.Sprintf("unknown job type: %s", job.Type)
		return
	}

	defer output.Close()

	if job.Key != "" {
		// enable encryption
		output, err = app.EncryptedWriter(job.Key, output)
		if err != nil {
			info.Status = BackupJobFailure
			info.Message = fmt.Sprintf("encryption failed: %v", err)
			return
		}
	}

	filecount := 0
	bytes, err := a.StreamZip(a.WorldDir, output, &app.StreamZipOpts{
		Callback: func(files, bytes int64, currentfile string) {
			info.Message = fmt.Sprintf("Copying file '%s' (progress: %d bytes, %d files)", currentfile, bytes, files)
			filecount++
		},
	})
	if err != nil {
		info.Status = BackupJobFailure
		info.Message = fmt.Sprintf("backup failed: %v", err)
		return
	}

	info.Message = fmt.Sprintf("Backup complete with %d bytes and %d files", bytes, filecount)
	info.Status = BackupJobSuccess

	a.CreateUILogEntry(&types.Log{
		Username: c.Username,
		Event:    "backup",
		Message:  info.Message,
	}, nil)
}

func (a *Api) CreateBackupJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	job := &CreateBackupJob{}
	err := json.NewDecoder(r.Body).Decode(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json decode error: %v", err))
		return
	}

	if job.ID == "" {
		job.ID = uuid.NewString()
	}

	info := &BackupJobInfo{
		Status: BackupJobRunning,
		ID:     job.ID,
	}
	backupjobs[job.ID] = info
	go backupJob(a.app, job, info, c)

	SendJson(w, info)
}

func (a *Api) GetBackupJobInfo(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	id := mux.Vars(r)["id"]
	info := backupjobs[id]
	if info == nil {
		SendError(w, 404, fmt.Errorf("backup job not found: %s", id))
		return
	}
	SendJson(w, info)
}
