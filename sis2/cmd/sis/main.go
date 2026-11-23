package main

import (
	"errors"
	"fmt"

	"sis/internal/task"
)

func main() {
	svc := task.NewService(task.TextAnalyzer{})
	svc.Start(3, 10) // 3 workers, queue of 10

	inputs := []string{
		"Hello Go world",
		"Go go GO",
		"line one\nline two",
		"привет мир",
	}

	var ids []int
	for _, in := range inputs {
		t, err := svc.CreateTask(in)
		if err != nil {
			fmt.Println("create:", err)
			continue
		}
		if err := svc.Submit(t.ID); err != nil {
			if errors.Is(err, task.ErrQueueFull) {
				fmt.Println("queue full, try later:", err)
			} else {
				fmt.Println("submit:", err)
			}
			continue
		}
		ids = append(ids, t.ID)
	}

	svc.Shutdown() // waits for every queued task to finish

	for _, id := range ids {
		t, err := svc.GetTaskByID(id)
		if err != nil {
			fmt.Println("get:", err)
			continue
		}
		fmt.Printf("task %d [%s]: words=%d runes=%d lines=%d unique=%d\n",
			t.ID, t.State, t.Result.Words, t.Result.Runes, t.Result.Lines, t.Result.UniqueWords)
	}

	if err := svc.Submit(1); err != nil {
		fmt.Println("after shutdown:", err)
	}
}
