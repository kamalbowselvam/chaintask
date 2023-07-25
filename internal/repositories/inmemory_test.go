package repositories

import (
	"testing"
)



func TestCreateTaskInMemory(t *testing.T) {
	generateRandomTask(t, testInMemoryStore)

}


func TestGetTaskInMemory(t *testing.T){
	GetTaskHelper(t, testInMemoryStore)
}

func TestGetTaskListInMemory(t *testing.T){
	GetTaskListHelper(t, testInMemoryStore)
}


func TestDeleteTaskInMemory(t *testing.T){
	DeleteTaskHelper(t, testInMemoryStore)
}

func TestUpdateTask(t *testing.T){
	UpdateTaskHelper(t, testInMemoryStore)
}