package domain

import (
	"time"
)

type Location [2]float64

type Project struct {
	Id                   int64     `json:"id"`
	Projectname          string    `json:"projectname"`
	CreatedOn            time.Time `json:"created_on"`
	CreatedBy            string    `json:"created_by"`
	Location             Location  `json:"location"`
	Address              string    `json:"address"`
	Responsible          string    `json:"responsible"`
	Client               string    `json:"client"`
	Tasks                []Task    `json:"tasks"`
	CompletionPercentage float64   `json:"completion_percentage"`
	Budget               float64   `json:"budget"`
	CompanyId            int64     `json:"company_id"`
}

func NewProject(projectname string, address string, location Location, responsible string, client string, user string) Project {

	return Project{
		Projectname: projectname,
		CreatedOn:   time.Now(),
		CreatedBy:   user,
		Location:    location,
		Address:     address,
		Responsible: responsible,
		Client:      client,
	}
}