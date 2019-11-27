package main

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Task struct {
	Name         string
	Priority     int
	Dependencies []*Task
}

func ParseFile(filename string) (n int64, m int64, tasks []*Task, err error) {
	// Read file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, 0, []*Task{}, err
	}

	// Parse data
	splits := strings.Split(string(file), "\n")

	// Parse header
	header := strings.Split(splits[0], " ")
	if len(header) < 2 {
		return 0, 0, []*Task{}, errors.New("invalid header")
	}
	n, err = strconv.ParseInt(header[0], 10, 64)
	if err != nil {
		return 0, 0, []*Task{}, err
	}
	m, err = strconv.ParseInt(header[1], 10, 64)
	if err != nil {
		return 0, 0, []*Task{}, err
	}

	for _, split := range splits[1:] {
		taskNames := strings.Split(split, "-")
		var prev *Task
		for _, taskName := range taskNames {
			task := FindTask(taskName, tasks)
			if task == nil {
				task = &Task{
					Name:         taskName,
					Priority:     0,
					Dependencies: []*Task{},
				}
				tasks = append(tasks, task)
			}
			if prev != nil {
				if FindTask(taskName, prev.Dependencies) == nil {
					prev.Dependencies = append(prev.Dependencies, task)
				}
			}
			prev = task
		}
	}

	return
}

func FindTask(name string, tasks []*Task) *Task {
	for _, task := range tasks {
		if task.Name == name {
			return task
		}
	}
	return nil
}

func ComputePriorities(tasks []*Task) error {
	return nil
}

func main() {
	n, m, tasks, err := ParseFile("tasks")
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}
	err = ComputePriorities(tasks)
	if err != nil {
		log.Fatalf("Failed to compute priorities: %v", err)
	}

	log.Printf("N = %v, M = %v", n, m)
}
