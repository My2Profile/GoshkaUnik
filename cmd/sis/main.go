package main

import (
	"fmt"
	"sis/internal/task"
)

func main() {
	taskService := task.NewService(task.TextAnalyzer{})

	taskService.CreateTask("Hello Go World")

	taskService.ProcessTask(1)

	result, err := taskService.GetTaskByID(1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(result)

}
