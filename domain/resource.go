package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Resource struct {
	Id           int64           `json:"id"`
	ResourceName string          `json:"resource_name"`
	Availed      decimal.Decimal `json:"availed"`
	Current      decimal.Decimal
	CreatedOn    time.Time
	CreatedBy    string
	UpdatedOn    time.Time
	UpdatedBy    string
}
