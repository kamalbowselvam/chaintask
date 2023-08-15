package domain

import "time"

type Task struct {
	Id        int64     `json:"id"`
	TaskName  string    `json:"taskname"`
	Budget    float64   `json:"budget"`
	CreatedOn time.Time `json:"createdOn"`
	CreatedBy string    `json:"createdBy"`
	UpdatedOn time.Time `json:"updatedOn"`
	UpdatedBy string    `json:"updatedBy"`
	Done      bool      `json:"done"`
}

func NewTask(name string, budget float64, user string) Task {

	t := Task{
		TaskName:  name,
		Budget:    budget,
		CreatedOn: time.Now(),
		CreatedBy: user,
		UpdatedOn: time.Now(),
		UpdatedBy: user,
		Done:      false,
	}
	return t
}

func (task *Task) GetTaskName() string {
	return task.TaskName
}

func (task *Task) GetBudget() float64 {
	return task.Budget
}

func (task *Task) IsDone() bool {
	return task.Done
}

func (task *Task) SetTaskDone(val bool) {
	task.Done = val
}
