package models

import "time"

type RedirectLog struct {
	Timestamp time.Time
	Ip        string
	Keyword   string
	UserAgent string
	Url       string
}
