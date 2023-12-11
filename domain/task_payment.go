package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type TaskPayment struct {
	Id        int64            `json:"id"`
	TaskId    int64            `json:"task_id"`
	Amount    decimal.Decimal  `json:"amount"`
	Source     string          `json:"source"`
	CreatedOn time.Time
	CreatedBy string
	UpdatedOn time.Time
	UpdatedBy string
}