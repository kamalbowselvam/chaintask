package domain

import "time"

type Payee struct {
	Id        int64  `json:"id"`
	PayeeName string `json:"payee_name"`
	Address   string `json:"address"`
	Current   float64
	CreatedOn time.Time
	CreatedBy string
	UpdatedOn time.Time
	UpdatedBy string
}
