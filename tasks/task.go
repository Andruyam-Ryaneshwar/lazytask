package task

import (
	"encoding/json"
	"fmt"
)

type StatusCode int

const (
	TODO StatusCode = iota
	INPROGRESS
	DONE
)

type Task struct {
	ID          int
	Title       string
	Description string
	StatusCode  int
	SubTasks    []*Task
}

func statusCodeToString(statusCode int) string {
	return [3]string{"TODO", "IN PROGRESS", "DONE"}[statusCode]
}

// TODO: use uuid ??
func GetNextID() int {
	return len(Tasks) + 1
}

var Tasks []*Task

// TODO: create a save function that after every CRUD operation, update the tasks in a json file
// TODO: add proper logs in CRUD operations

func AddTask(task *Task) {
	task.ID = GetNextID()
	Tasks = append(Tasks, task)
	fmt.Println("Task created successfully!")
}

// TODO: if sticking with the increment id thing, we can probably use binary search here
func UpdateTask(updatedTask *Task) {
	for i, task := range Tasks {
		if task.ID == updatedTask.ID {
			Tasks[i] = updatedTask
			break
		}
	}
	fmt.Println("Task updated successfully!")
}

func GetTasks() []*Task {
	return Tasks
}

func DeleteTask(taskId int) {
	for i, task := range Tasks {
		if task.ID == taskId {
			Tasks = append(Tasks[:i], Tasks[i+1:]...)
		}
	}
	fmt.Println("Task deleted successfully")
}
