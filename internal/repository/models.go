package repository

import "time"

type ClickStat struct {
	Timestamp time.Time `json:"timestamp"`
	BannerID  int       `json:"banner_id"`
	Count     int       `json:"count"`
}
