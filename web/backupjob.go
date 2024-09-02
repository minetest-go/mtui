package web

import (
	"encoding/json"
	"fmt"
	"mtui/app"
	"mtui/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/sftp"
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
	BackupJobTypeSCP BackupJobType = "scp"
)

type CreateBackupJob struct {
	Type     BackupJobType `json:"type"`
	Host     string        `json:"host"`
	Port     int           `json:"port"`
	Filename string        `json:"filename"`
	Username string        `json:"username"`
	Password string        `json:"password"`
}

var backupjobs = map[string]*BackupJobInfo{}

func backupJob(a *app.App, job *CreateBackupJob, info *BackupJobInfo) {
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

	sc, err := sftp.NewClient(client)
	if err != nil {
		info.Status = BackupJobFailure
		info.Message = fmt.Sprintf("sftp open failed: %v", err)
		return
	}
	defer sc.Close()

	file, err := sc.Create(job.Filename)
	if err != nil {
		info.Status = BackupJobFailure
		info.Message = fmt.Sprintf("sftp create failed: %v", err)
		return
	}
	defer file.Close()

	filecount := 0
	bytes, err := a.StreamZip(a.WorldDir, file, &app.StreamZipOpts{
		Callback: func(files, bytes int64, currentfile string) {
			info.Message = fmt.Sprintf("Copying file '%s' (progress: %d bytes, %d files)", currentfile, bytes, files)
			filecount++
		},
	})
	if err != nil {
		info.Status = BackupJobFailure
		info.Message = fmt.Sprintf("sftp create failed: %v", err)
		return
	}

	info.Message = fmt.Sprintf("Backup complete with %d bytes and %d files", bytes, filecount)
	info.Status = BackupJobSuccess

	a.CreateUILogEntry(&types.Log{
		Event:   "backup",
		Message: info.Message,
	}, nil)
}

func (a *Api) CreateBackupJob(w http.ResponseWriter, r *http.Request) {
	job := &CreateBackupJob{}
	err := json.NewDecoder(r.Body).Decode(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json decode error: %v", err))
		return
	}

	id := uuid.NewString()
	info := &BackupJobInfo{
		Status: BackupJobRunning,
		ID:     id,
	}
	backupjobs[id] = info
	go backupJob(a.app, job, info)

	SendJson(w, info)
}

func (a *Api) GetBackupJobInfo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	info := backupjobs[id]
	if info == nil {
		SendError(w, 404, fmt.Errorf("backup job not found: %s", id))
		return
	}
	SendJson(w, info)
}
