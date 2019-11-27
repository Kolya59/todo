package main

import (
	"errors"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Task struct {
	Name         string
	Priority     int
	Dependencies []*Task
	Iteration    int
	Executor     int64
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
	k := 1
	for {
		sort.Slice(tasks, func(i, j int) bool {
			is, js := 0, 0
			for _, iDep := range tasks[i].Dependencies {
				is += iDep.Priority
			}
			for _, jDep := range tasks[j].Dependencies {
				js += jDep.Priority
			}
			return is < js
		})
		changesFlag := false
		for _, task := range tasks {
			if task.Priority == 0 {
				allowedToResolving := true
				for _, dep := range task.Dependencies {
					allowedToResolving = allowedToResolving && dep.Priority != 0
				}
				if allowedToResolving {
					changesFlag = true
					task.Priority = k
					k++
				}
			}
		}
		if !changesFlag {
			break
		}
	}

	isResolved := true
	for _, task := range tasks {
		isResolved = isResolved && task.Priority != 0
	}
	if !isResolved {
		return errors.New("probably, frequency contains cycle dependencies")
	}

	return nil
}

func SetResolvers(tasks []*Task, count int64) {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority < tasks[j].Priority
	})

	l := len(tasks)
	taskIdx := 0
	for i := 0; taskIdx < l; i++ {
		currExecutor := count
		for currExecutor > 0 && taskIdx < l {
			tasks[taskIdx].Executor = currExecutor
			tasks[taskIdx].Iteration = i + 1
			currExecutor--
			taskIdx++
		}
	}
}

func main() {
	n, m, tasks, err := ParseFile("tasks")
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}
	log.Printf("N = %v, M = %v", n, m)

	err = ComputePriorities(tasks)
	if err != nil {
		log.Fatalf("Failed to compute priorities: %v", err)
	}

	SetResolvers(tasks, m)

	for _, task := range tasks {
		log.Printf("Task: %v with priority: %v, executor: %v and iteration: %v", task.Name, task.Priority, task.Executor, task.Iteration)
	}

}
