package domain

import "time"

type Company struct {
	Id                   int64     `json:"id"`
	CompanyName          string    `json:"company_name"`
	Address              string    `json:"address"`
	CreatedOn            time.Time 
	CreatedBy            string    
	UpdatedOn            time.Time 
	UpdatedBy            string   
}
