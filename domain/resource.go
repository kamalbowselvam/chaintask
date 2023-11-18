package domain

import "time"

type Resource struct {
	Id           int64   `json:"id"`
	ResourceName string  `json:"resource_name"`
	Availed      float64 `json:"availed"`
	Current      float64
	CreatedOn    time.Time
	CreatedBy    string
	UpdatedOn    time.Time
	UpdatedBy    string
}
