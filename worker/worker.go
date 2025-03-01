package worker

import (
	"net/http"
	"time"
)

func NewWorkerHttpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

type Worker interface {
	Run() error
}

var Workers []Worker
