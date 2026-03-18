package httputil

import (
	"net"
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     30 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: 3 * time.Second,
		}).DialContext,
	},
}
