package domain

import (
	"time"

	"github.com/kamalbowselvam/chaintask/customtypes"
)

type Location [2]float64

type Project struct {
	Id            int64     `json:"id"`
	Projectname   string    `json:"projectname"`
	CreatedOn     time.Time `json:"createdOn"`
	CreatedBy     string    `json:"createdBy"`
	Location      Location  `json:"location"`
	LocationPoint customtypes.Point
	Address       string `json:"address"`
	Responsible   string `json:"responsible"`
	Client        string `json:"client"`
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
