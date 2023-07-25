package repositories

import (
	"testing"
)


func TestCreateTaskPersistence(t *testing.T) {
	generateRandomTask(t, testPersistenceStore)

}


func TestGetTaskPersistence(t *testing.T){
	GetTaskHelper(t, testPersistenceStore)
}

func TestGetTaskListPersistence(t *testing.T){
	GetTaskListHelper(t, testPersistenceStore)
}


func TestDeleteTaskPersistence(t *testing.T){
	DeleteTaskHelper(t, testPersistenceStore)
}


func TestUpdateTaskPersistence(t *testing.T) {
	UpdateTaskHelper(t, testPersistenceStore)
}