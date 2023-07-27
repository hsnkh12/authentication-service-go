package http

import (
	"net/http"
	"time"
)

func NewServer(address string) *http.Server {

	return &http.Server{
		Addr:           address,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
