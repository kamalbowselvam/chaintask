package models

import "time"

type Task struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Budget     float64   `json:"budget"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  string    `json:"createdBy"`
	ValidateBy string    `json:"vaidatedBy"`
	ValidateOn time.Time `json:"validatedOn"`
	UpdatedOn  time.Time `json:"updatedOn"`
	UpdatedBy  time.Time `json:"updatedBy"`
}
