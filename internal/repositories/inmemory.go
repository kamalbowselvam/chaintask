package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/kamalbowselvam/chaintask/internal/core/domain"
)


type InMemoryStorage struct {

	keyvaluestore map[int64]domain.Task
	entries int64 
}


func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		keyvaluestore: map[int64]domain.Task{},
		entries: 1,
	}
}



func (repo *InMemoryStorage) SaveTask(ctx context.Context, task domain.Task) (domain.Task, error){
	task.Id = repo.entries
	repo.keyvaluestore[repo.entries] = task
	repo.entries += 1
	var err error
	return task, err 
}


func (repo *InMemoryStorage) GetTask(ctx context.Context, id int64) (domain.Task,error){
	task, ok := repo.keyvaluestore[id]

	var err error

	if !ok {
		err = errors.New("key not found")
		fmt.Println("Here")
		return domain.Task{}, err
	}

	return task, err
}


func (repo *InMemoryStorage) GetTaskList( ctx context.Context, ids []int64) ([]domain.Task, error){
	var tasks []domain.Task
	var err error

	for i, s := range ids {
		task, ok := repo.keyvaluestore[s]
		if !ok {
			fmt.Println("Key not found")
		} else {
			tasks = append(tasks, task)
		}
		
		fmt.Println(i, s)
	}
	return tasks, err
}


func (repo *InMemoryStorage) DeleteTask(ctx context.Context, id int64) error {
	_, ok := repo.keyvaluestore[id]

	var err error

	if ok{
		delete(repo.keyvaluestore,id)
	} else {
		err = errors.New("key not found")
	}
	return err
}


func (repo *InMemoryStorage) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {

	id := task.Id
	var err error

	if id == 0 {
		err := errors.New("id is nil in the task")
		return domain.Task{}, err
	}
	repo.keyvaluestore[id] = task

	return task, err

}

