package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kamalbowselvam/chaintask/domain"
)

type InMemoryStorage struct {
	taskstore map[int64]domain.Task
	userstore map[string]domain.User
	entries   int64
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		taskstore: map[int64]domain.Task{},
		userstore: map[string]domain.User{},
		entries:   1,
	}
}

func (repo *InMemoryStorage) GetUser(ctx context.Context, username string) (domain.User, error){

	var err error
	return repo.userstore[username], err
}

func (repo *InMemoryStorage) CreateUser(ctx context.Context, arg CreateUserParams) (domain.User, error) {

	user := domain.User{
		Username: arg.Username,
		HashedPassword: arg.HashedPassword,
		FullName: arg.FullName,
		Email: arg.Email,
		CreatedAt: time.Now(),
	}
	repo.userstore[user.Username] = user
	var err error
	return user, err
}

func (repo *InMemoryStorage) CreateTask(ctx context.Context, arg CreateTaskParams) (domain.Task, error) {

	task := domain.Task{
		TaskName: arg.TaskName,
		Budget: arg.Budget,
		CreatedBy: arg.CreatedBy,
		CreatedOn: time.Now(),
		UpdatedBy: arg.CreatedBy,
		UpdatedOn: time.Now(),
		Done: false,
	}

	task.Id = repo.entries
	repo.taskstore[repo.entries] = task
	repo.entries += 1
	var err error
	return task, err
}

func (repo *InMemoryStorage) GetTask(ctx context.Context, arg GetTaskParams) (domain.Task, error) {
	task, ok := repo.taskstore[arg.Id]

	var err error

	if !ok {
		err = errors.New("key not found")
		fmt.Println("Here")
		return domain.Task{}, err
	}

	return task, err
}

func (repo *InMemoryStorage) GetTaskList(ctx context.Context, ids []int64) ([]domain.Task, error) {
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

	if ok {
		delete(repo.taskstore, id)
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
