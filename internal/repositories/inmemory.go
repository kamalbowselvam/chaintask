package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kamalbowselvam/chaintask/internal/core/domain"
)


type InMemoryStorage struct {

	taskstore map[int64]domain.Task
	userstore map[string]domain.User
	entries int64 
}


func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		taskstore: map[int64]domain.Task{},
		userstore: map[string]domain.User{},
		entries: 1,
	}
}


func (repo *InMemoryStorage) CreateUser(ctx context.Context, user domain.User) (domain.UserDetail, error){
	repo.userstore[user.Username] = user
	var err error

	var userdetail domain.UserDetail

	userdetail.Username = user.Username
	userdetail.Email = user.Email
	userdetail.CreatedAt = time.Now()
	userdetail.FullName = user.FullName

	return userdetail, err

}



func (repo *InMemoryStorage) SaveTask(ctx context.Context, task domain.Task) (domain.Task, error){
	task.Id = repo.entries
	repo.taskstore[repo.entries] = task
	repo.entries += 1
	var err error
	return task, err 
}


func (repo *InMemoryStorage) GetTask(ctx context.Context, id int64) (domain.Task,error){
	task, ok := repo.taskstore[id]

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
		task, ok := repo.taskstore[s]
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
	_, ok := repo.taskstore[id]

	var err error

	if ok{
		delete(repo.taskstore,id)
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
	repo.taskstore[id] = task

	return task, err

}

