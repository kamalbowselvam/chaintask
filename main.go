package main

import (
	"fmt"
	"time"

	"github.com/kamalbowselvam/chaintask/models"
)

func main(){


	fmt.Println("Hello World Task")

	task := &models.Task{
		ID: "1",
		Name: "Kamal",
		Budget: 1000.0,
		CreatedAt:  time.Now(),
		CreatedBy: "Kamal",
		ValidateBy: "Menelik" ,   
		ValidateOn:  time.Now() ,
		UpdatedOn:   time.Now() ,
		UpdatedBy:   time.Now() ,
	}



	fmt.Println(task.Name)



	return 

}

