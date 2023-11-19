package domain

import (
	"time"
)

type Project struct {
	Id                   int64     `json:"id"`
	Projectname          string    `json:"project_name"`
	CreatedOn            time.Time `json:"created_on"`
	CreatedBy            string    `json:"created_by"`
	Longitude             float64  `json:"longitude"`
	Latitude             float64  `json:"latitude"`
	Address              string    `json:"address"`
	Responsible          string    `json:"responsible"`
	Client               string    `json:"client"`
	Tasks                []Task    `json:"tasks"`
	CompletionPercentage float64   `json:"completion_percentage"`
	Budget               float64   `json:"budget"`
	CompanyId            int64     `json:"company_id"`
}

func NewProject(projectname string, address string, longitude float64, latitude float64, responsible string, client string, user string) Project {

	return Project{
		Projectname: projectname,
		CreatedOn:   time.Now(),
		CreatedBy:   user,
		Longitude:    longitude,
		Latitude: latitude,
		Address:     address,
		Responsible: responsible,
		Client:      client,
	}
}