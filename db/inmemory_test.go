package db

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

func TestUpdateTaskInMemory(t *testing.T){
	UpdateTaskHelper(t, testInMemoryStore)
}

func TestGetUserInMemory(t *testing.T){
	GetUserHelper(t, testInMemoryStore)
}