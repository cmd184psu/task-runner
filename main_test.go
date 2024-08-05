package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// func TestExecuteSSHCommand(t *testing.T) {
// 	expected := "echo Task 1"
// 	output, err := executeSSHCommand(expected)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}
// 	if output != expected {
// 		t.Fatalf("Expected %s, got %s", expected, output)
// 	}
// }

func TestWorker(t *testing.T) {
	task := Task{ID: 1, Command: "echo Test"}
	taskQueue = make(chan Task, 1)
	results = sync.Map{}
	taskQueue <- task
	close(taskQueue)

	wg.Add(1)
	go worker()
	wg.Wait()

	value, ok := results.Load(task.ID)
	if !ok {
		t.Fatalf("Expected task result to be stored")
	}
	storedTask, ok := value.(Task)
	if !ok {
		t.Fatalf("Expected stored value to be of type Task")
	}
	if storedTask.ID != task.ID || storedTask.Command != task.Command {
		t.Fatalf("Expected task %v, got %v", task, storedTask)
	}
}

func TestHandleStatus(t *testing.T) {
	task := Task{ID: 1, Command: "echo Test", StartTime: time.Now()}
	results.Store(task.ID, task)

	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, status)
	}

	expected := `[{"id":1,"command":"echo Test","start_time":`
	if rr.Body.String()[:len(expected)] != expected {
		t.Fatalf("Expected response to start with %s, got %s", expected, rr.Body.String())
	}
}
