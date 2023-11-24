package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Task struct {
	Id        int64           `json:"id"`
	TaskName  string          `json:"taskname"`
	Budget    decimal.Decimal `json:"budget"`
	CreatedOn time.Time
	CreatedBy string
	UpdatedOn time.Time
	UpdatedBy string
	Done      bool  `json:"done"`
	TaskOrder int64 `json:"task_order"`
	ProjectId int64 `json:"project_id"`
	Version   int64 `json:"version"`
	CompanyId int64 `json:"company_id"`
}

func NewTask(name string, budget decimal.Decimal, user string, taskorder int64, projectid int64) Task {

	t := Task{
		TaskName:  name,
		Budget:    budget,
		CreatedOn: time.Now(),
		CreatedBy: user,
		UpdatedOn: time.Now(),
		UpdatedBy: user,
		Done:      false,
		TaskOrder: taskorder,
		ProjectId: projectid,
	}
	return t
}

func (task *Task) GetTaskName() string {
	return task.TaskName
}

func (task *Task) GetBudget() decimal.Decimal {
	return task.Budget
}

func (task *Task) IsDone() bool {
	return task.Done
}

func (task *Task) SetTaskDone(val bool) {
	task.Done = val
}
