package config

import "time"

const (
	Protocol = "http"
	Host     = "localhost"
	Port     = 8000
)
const (
	SessionID           = "sessionID"
	CookieExpireSeconds = 120
	RequestTimeout      = 3 * time.Second
)
