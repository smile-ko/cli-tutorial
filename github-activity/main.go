package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		Commits []any  `json:"commits"`
	} `json:"payload"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a GitHub username.")
		return
	}

	userName := os.Args[1]

	url := fmt.Sprintf("https://api.github.com/users/%s/events", userName)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Printf("Recent events for user %s:\n", userName)
	fmt.Println("=============================================")
	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			fmt.Printf("- Pushed %d commits to %s\n", len(event.Payload.Commits), event.Repo.Name)
		case "IssuesEvent":
			fmt.Printf("- %s a new issue in %s\n", capitalize(event.Payload.Action), event.Repo.Name)
		case "WatchEvent":
			fmt.Printf("- Starred %s\n", event.Repo.Name)
		default:
		}
	}
	fmt.Println("=============================================")
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
