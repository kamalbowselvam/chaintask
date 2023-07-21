package models

import "time"

type Task struct {
	Id        int64    `json:"id"`
	Name      string    `json:"name"`
	Budget    float64   `json:"budget"`
	CreatedOn time.Time `json:"createdOn"`
	CreatedBy string    `json:"createdBy"`
	UpdatedOn time.Time `json:"updatedOn"`
	UpdatedBy string    `json:"updatedBy"`
	Done      bool      `json:"done"`
}
