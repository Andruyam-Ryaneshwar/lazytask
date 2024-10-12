package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID          int
	Title       string
	Description string
	StatusCode  int
	SubTasks    []*Task
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var tasks []*Task

func init() {
	loadTasks()
}

// Private functions
func getNextID() int {
	return len(tasks) + 1
}

func saveTasks() {
	data, err := json.Marshal(tasks)

	if err != nil {
		fmt.Println("Unable to save tasks!", err)
		return
	}

	err = os.WriteFile("tasks.json", data, os.ModePerm)

	if err != nil {
		fmt.Println("Unable to save file!")
	}
}

func loadTasks() {
	data, err := os.ReadFile("tasks.json")

	if err != nil {
		fmt.Println("Unable to read the file!")
		return
	}

	err = json.Unmarshal(data, &tasks)

	if err != nil {
		fmt.Println("Unable to parse the tasks!", err)
	}
}

// Public functions
func StatusCodeToString(statusCode int) string {
	statuses := [3]string{"TODO", "IN PROGRESS", "DONE"}
	if statusCode >= 0 && statusCode < len(statuses) {
		return statuses[statusCode]
	}
	return "UKNOWN"
}

func AddTask(task *Task) {
	task.ID = getNextID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	tasks = append(tasks, task)
	fmt.Println("Task created successfully!")
	saveTasks()
}

func UpdateTask(updatedTask *Task) {
	for i, task := range tasks {
		if task.ID == updatedTask.ID {
			updatedTask.UpdatedAt = time.Now()
			tasks[i] = updatedTask
			fmt.Println("Task updated successfully!")
			saveTasks()
			return
		}
	}
	fmt.Println("Task not found!")
}

func GetTasks() []*Task {
	return tasks
}

func GetTaskByID(id int) *Task {
	for _, task := range tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func DeleteTask(taskId int) {
	for i, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Println("Task deleted successfully")
			saveTasks()
			return
		}
	}
	fmt.Println("Task not found!")
}
