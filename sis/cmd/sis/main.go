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

	resultData := result.Result
	fmt.Printf("Task ID: %d\n", result.ID)
	fmt.Printf("State: %s\n", result.State)
	fmt.Printf("Words: %d\n", resultData.Words)
	fmt.Printf("Runes: %d\n", resultData.Runes)
	fmt.Printf("Lines: %d\n", resultData.Lines)
	fmt.Printf("Unique Words: %d\n", resultData.UniqueWords)

}
