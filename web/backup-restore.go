package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mtui/app"
	"mtui/types"
	"net/http"
	"sync/atomic"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type BackupRestoreType string

const (
	BackupJob  BackupRestoreType = "backup"
	RestoreJob BackupRestoreType = "restore"
)

type BackupRestoreJobState string

const (
	BackupRestoreJobRunning BackupRestoreJobState = "running"
	BackupRestoreJobSuccess BackupRestoreJobState = "success"
	BackupRestoreJobFailure BackupRestoreJobState = "failure"
)

// new job
type CreateBackupRestoreJob struct {
	Type BackupRestoreType `json:"type"`

	Endpoint  string `json:"endpoint"`
	KeyID     string `json:"key_id"`
	AccessKey string `json:"access_key"`
	Bucket    string `json:"bucket"`

	FileKey  string `json:"file_key"`
	Filename string `json:"filename"`
}

// current job info
type BackupRestoreInfo struct {
	Type            BackupRestoreType     `json:"type"`
	ProgressPercent float64               `json:"progress_percent"`
	Message         string                `json:"message"`
	State           BackupRestoreJobState `json:"state"`
}

var backupRestoreInfo = atomic.Pointer[BackupRestoreInfo]{}

func getS3Client(job *CreateBackupRestoreJob) (*minio.Client, error) {
	return minio.New(job.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(job.KeyID, job.AccessKey, ""),
		Secure: true,
	})
}

func (a *Api) startBackupJob(job *CreateBackupRestoreJob) error {
	// calculate estimated directory size
	estimated_size, err := a.getDirectorySize(a.app.WorldDir)
	if err != nil {
		return fmt.Errorf("directory size estimation: %v", err)
	}

	s3, err := getS3Client(job)
	if err != nil {
		return fmt.Errorf("get s3 client error: %v", err)
	}

	pr, pw := io.Pipe()

	var output io.Writer
	if job.FileKey != "" {
		output, err = app.EncryptedWriter(job.FileKey, pw)
		if err != nil {
			return fmt.Errorf("encryption setup: %v", err)
		}
	}

	// write zip file to output
	go func() {
		_, err := a.app.StreamZip(a.app.WorldDir, output, &app.StreamZipOpts{
			Callback: func(files, bytes int64, currentfile string) {
				progress_percent := float64(bytes) / float64(estimated_size) * 100

				backupRestoreInfo.Store(&BackupRestoreInfo{
					Type:            job.Type,
					ProgressPercent: progress_percent,
					State:           BackupRestoreJobRunning,
					Message:         fmt.Sprintf("Copying '%s' (%d / %d bytes)", currentfile, bytes, estimated_size),
				})
			},
		})
		if err != nil {
			pw.CloseWithError(fmt.Errorf("stream zip error: %v", err))
		} else {
			pw.Close()
		}
	}()

	// read file content from input
	_, err = s3.PutObject(context.Background(), job.Bucket, job.Filename, pr, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("s3 upload error: %v", err)
	}

	return nil
}

func (a *Api) startRestoreJob(job *CreateBackupRestoreJob) error {
	// TODO
	return nil
}

// api

func (a *Api) CreateBackupRestoreJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	info := backupRestoreInfo.Load()
	if info != nil && info.State == BackupRestoreJobRunning {
		SendError(w, 405, fmt.Errorf("job already running"))
		return
	}

	job := &CreateBackupRestoreJob{}
	err := json.NewDecoder(r.Body).Decode(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json error: %v", err))
		return
	}

	var job_fn func(*CreateBackupRestoreJob) error

	switch job.Type {
	case BackupJob:
		job_fn = a.startBackupJob
	case RestoreJob:
		job_fn = a.startRestoreJob
	default:
		SendError(w, 500, fmt.Errorf("unknown type: '%s'", job.Type))
		return
	}

	// start job
	go func() {
		err := job_fn(job)
		if err != nil {
			backupRestoreInfo.Store(&BackupRestoreInfo{
				Type:    job.Type,
				Message: fmt.Sprintf("job failed: %v", err),
				State:   BackupRestoreJobFailure,
			})
		} else {
			backupRestoreInfo.Store(&BackupRestoreInfo{
				Type:            job.Type,
				Message:         "done",
				ProgressPercent: 100,
				State:           BackupRestoreJobSuccess,
			})
		}
	}()

	// job info
	info = &BackupRestoreInfo{
		Type:    job.Type,
		Message: "Starting",
		State:   BackupRestoreJobRunning,
	}
	backupRestoreInfo.Store(info)

	Send(w, info, nil)
}

func (a *Api) GetBackupRestoreJobInfo(w http.ResponseWriter, r *http.Request) {
	info := backupRestoreInfo.Load()
	if info == nil {
		SendError(w, 404, fmt.Errorf("no job started"))
		return
	}
	Send(w, info, nil)
}
