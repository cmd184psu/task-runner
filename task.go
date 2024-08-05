package main

import (
	"sync"
	"time"

	"github.com/cmd184psu/alfredo"
)

// type Task struct {
// 	ID        int           `json:"id"`
// 	Command   string        `json:"command"`
// 	StartTime time.Time     `json:"start_time"`
// 	Duration  time.Duration `json:"duration"`
// 	Output    string        `json:"output"`
// 	Error     string        `json:"error"`
// 	IsRunning bool          `json:"is_running"`
// 	ssh       alfredo.SSHStruct
// }

type Task struct {
	ID        int           `json:"id"`
	Command   string        `json:"command"`
	StartTime time.Time     `json:"start_time"`
	Duration  time.Duration `json:"duration"`
	Output    string        `json:"output"`
	IsRunning bool          `json:"is_running"`
	Error     string        `json:"error"`
	ssh       alfredo.SSHStruct
	HasRun    bool `json:"has_run"`
}

var (
	taskQueue   = make(chan Task, 100)
	queuedTasks = []Task{} // Separate slice to track queued tasks
	results     = sync.Map{}
	taskCounter = 0
	wg          sync.WaitGroup
	stop        = make(chan struct{})
	ssh         = alfredo.SSHStruct{Key: alfredo.ExpandTilde("~/.ssh/homelab_rsa"), Host: "192.168.1.10", User: "cdelezenski"}
)

func Enqueue(s string) {
	if len(s) == 0 {
		panic("attempted to enqueue command of zero length")
	}

	task := Task{
		ID:      taskCounter,
		Command: s,
		ssh:     ssh,
	}
	taskCounter++

	taskQueue <- task
	queuedTasks = append(queuedTasks, task)
	results.Store(task.ID, task)
}
