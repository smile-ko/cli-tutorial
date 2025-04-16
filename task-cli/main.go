package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	fileName = "tasks.json"

	// types of tasks
	todo       = "TODO"
	inProgress = "IN_PROGRESS"
	done       = "DONE"

	// format for displaying time
	displayTime = "02/01/2006"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli [add|update|delete|mark-in-progress|mark-done|list]")
		return
	}

	command := os.Args[1]
	tasks := LoadTasks()

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli add <task>")
			return
		}

		id := generateID()
		title := os.Args[2]
		now := time.Now()
		task := Task{ID: id, Title: title, Status: todo, CreatedAt: now, UpdatedAt: now}
		_ = append(tasks, task)
		saveTasks(tasks)
		fmt.Printf("Task added: %s. [%s] %s (createdAt: %s)\n", id, task.Status, task.Title, task.CreatedAt.Format(displayTime))
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: task-cli update <task-id> <new-task>")
			return
		}
		id := os.Args[2]
		title := os.Args[3]
		for i, t := range tasks {
			if t.ID == id {
				tasks[i].Title = title
			}
		}
		saveTasks(tasks)
		fmt.Printf("Task updated: %s. [%s] %s (updatedAt: %s)\n", id, tasks[0].Status, title, tasks[0].UpdatedAt.Format(displayTime))
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli delete <task-id>")
			return
		}

		id := os.Args[2]
		newTasks := []Task{}
		for _, t := range tasks {
			if t.ID != id {
				newTasks = append(newTasks, t)
			}
		}
		saveTasks(newTasks)
		fmt.Printf("Task deleted: %s\n", id)
	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli mark-in-progress <task-id>")
			return
		}

		id := os.Args[2]
		for i, t := range tasks {
			if t.ID == id {
				tasks[i].Status = inProgress
			}
		}
		saveTasks(tasks)
		fmt.Printf("Task marked in progress: %s. [%s] %s (updatedAt: %s)\n", id, tasks[0].Status, tasks[0].Title, tasks[0].UpdatedAt.Format(displayTime))
	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli mark-in-progress <task-id>")
			return
		}

		id := os.Args[2]
		for i, t := range tasks {
			if t.ID == id {
				tasks[i].Status = inProgress
			}
		}
		saveTasks(tasks)
		fmt.Printf("Task marked in done: %s. [%s] %s (updatedAt: %s)\n", id, tasks[0].Status, tasks[0].Title, tasks[0].UpdatedAt.Format(displayTime))
	case "list":
		if len(os.Args) == 2 {
			for _, t := range tasks {
				fmt.Printf("%s. [%s] %s %s\n", t.ID, t.Status, t.Title, t.CreatedAt.Format(displayTime))
			}
		} else {
			filter := os.Args[2]
			for _, t := range tasks {
				if t.Status == filter {
					fmt.Printf("%s. [%s] %s %s\n", t.ID, t.Status, t.Title, t.CreatedAt.Format(displayTime))
				}
			}
		}
	default:
		fmt.Println("Usage: task-cli [add|update|delete|mark-in-progress|mark-done|list]")
	}
}

func LoadTasks() []Task {
	var tasks []Task
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks // Return empty slice if file does not exist
		}
		return nil
	}
	defer file.Close()

	stat, _ := file.Stat()
	if stat.Size() == 0 {
		return tasks // Return empty slice if file is empty
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return tasks
	}

	return tasks
}

func saveTasks(tasks []Task) {
	data, _ := json.MarshalIndent(tasks, "", "  ")
	os.WriteFile(fileName, data, 0644)
}

func generateID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic(err)
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16],
	)
}

func AddTask(title string) {
}

func updateTask(id string, title string) {
	tasks := LoadTasks()
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Title = title
			tasks[i].UpdatedAt = time.Now()
		}
	}
	saveTasks(tasks)
	fmt.Printf("Task updated: %s. [%s] %s (updatedAt: %s)\n", id, tasks[0].Status, title, tasks[0].UpdatedAt.Format(displayTime))
}
