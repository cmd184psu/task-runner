package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// func handleCompleted(w http.ResponseWriter, r *http.Request) {
// 	var completedTasks []Task
// 	//now := time.Now()
// 	results.Range(func(key, value interface{}) bool {
// 		task := value.(Task)
// 		if !task.IsRunning {
// 			//task.Duration = now.Sub(task.StartTime)
// 			completedTasks = append(completedTasks, task)
// 		}
// 		return true
// 	})

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(completedTasks); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func handleRunning(w http.ResponseWriter, r *http.Request) {
// 	var runningTasks []Task
// 	now := time.Now()
// 	results.Range(func(key, value interface{}) bool {
// 		task := value.(Task)
// 		if task.IsRunning {
// 			task.Duration = now.Sub(task.StartTime)
// 			runningTasks = append(runningTasks, task)
// 		}
// 		return true
// 	})

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(runningTasks); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func handleAddTask(w http.ResponseWriter, r *http.Request) {
	var newTask struct {
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.EqualFold("init", newTask.Command) {
		log.Println("== added new command init ==")
		//initTasks(10)

		// func initTasks(numTasks int) {
		var r int
		for i := 0; i < 20; i++ {
			r = rand.Intn(20) + 10
			Enqueue(fmt.Sprintf("sleep %d; echo Task %d", r, i))
		}

	} else {
		Enqueue(newTask.Command)
		// 	if len(newTask.Command) == 0 {
		// 		panic("attempted to run blank command")
		// 	}

		// 	task := Task{
		// 		ID:      taskCounter,
		// 		Command: newTask.Command,
		// 		ssh:     ssh,
		// 	}
		// 	taskCounter++

		// 	taskQueue <- task
		// 	results.Store(task.ID, task)
	}
	log.Printf("size of queue: %d", len(taskQueue))
	w.WriteHeader(http.StatusCreated)
}

// func handleQueued(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(queuedTasks); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func handleQueued(w http.ResponseWriter, r *http.Request) {
// 	var queuedTasks []Task
// 	taskQueueLen := len(taskQueue)
// 	for i := 0; i < taskQueueLen; i++ {
// 		task := <-taskQueue
// 		queuedTasks = append(queuedTasks, task)
// 		taskQueue <- task
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(queuedTasks); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func handleStop(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is shutting down..."))
	close(stop)
	go func() {
		wg.Wait()
		saveResults()
		os.Exit(0)
	}()
}

// func handleStatus(w http.ResponseWriter, r *http.Request) {
// 	var queued, running, completed []Task
// 	now := time.Now()

// 	// Get queued tasks
// 	for _, task := range queuedTasks {
// 		queued = append(queued, task)
// 	}

// 	// Get running and completed tasks
// 	results.Range(func(key, value interface{}) bool {
// 		task := value.(Task)
// 		if task.IsRunning {
// 			task.Duration = now.Sub(task.StartTime)
// 			running = append(running, task)
// 		} else {
// 			completed = append(completed, task)
// 		}
// 		return true
// 	})

// 	total := len(queued) + len(running) + len(completed)
// 	progress := float64(len(completed)) / float64(total) * 100

// 	response := struct {
// 		Queued    []Task  `json:"queued"`
// 		Running   []Task  `json:"running"`
// 		Completed []Task  `json:"completed"`
// 		Progress  float64 `json:"progress"`
// 	}{
// 		Queued:    queued,
// 		Running:   running,
// 		Completed: completed,
// 		Progress:  progress,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func handleStatus(w http.ResponseWriter, r *http.Request) {
	queued := []Task{}
	running := []Task{}
	completed := []Task{}
	now := time.Now()

	// Get queued tasks
	//	queued = append(queued, queuedTasks...)
	queued = append(queued, queuedTasks...)

	log.Printf("handleStatus::size of queue is %d", len(queued))

	// Get running and completed tasks
	results.Range(func(key, value interface{}) bool {
		task := value.(Task)
		if task.IsRunning {
			task.Duration = now.Sub(task.StartTime)
			running = append(running, task)
		} else if task.HasRun {
			completed = append(completed, task)
		}
		return true
	})

	total := len(queued) + len(running) + len(completed)
	var progress float64
	if total > 0 {
		progress = float64(len(completed)) / float64(total) * 100
	} else {
		progress = 0
	}

	response := struct {
		Queued    []Task  `json:"queued"`
		Running   []Task  `json:"running"`
		Completed []Task  `json:"completed"`
		Progress  float64 `json:"progress"`
	}{
		Queued:    queued,
		Running:   running,
		Completed: completed,
		Progress:  progress,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
