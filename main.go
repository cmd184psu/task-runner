package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sevlyar/go-daemon"
)

// func main() {
// 	alfredo.VerbosePrintln("hello world")

// 	var s alfredo.SSHStruct

// 	s.Host = "192.168.1.10"
// 	s.User = "cdelezenski"
// 	s.Key = alfredo.ExpandTilde("~/.ssh/homelab_rsa")

// 	// Initialize tasks
// 	initTasks(10, s)

// 	// Start worker pool
// 	startWorkers(10)

// 	// Start the HTTP server
// 	http.HandleFunc("/completed", handleCompleted)
// 	http.HandleFunc("/running", handleRunning)
// 	http.HandleFunc("/queued", handleQueued)
// 	http.HandleFunc("/add", handleAddTask)
// 	// Serve static files
// 	http.Handle("/", http.FileServer(http.Dir("./static")))
// 	// Channel to receive OS signals

// 	go func() {
// 		log.Println("Starting HTTP server on :8080")
// 		if err := http.ListenAndServe(":8080", nil); err != nil {
// 			log.Fatalf("Failed to start HTTP server: %v", err)
// 		}
// 	}()

// 	waitForCompletion()
// 	saveResults("results.json")
// 	log.Println("server terminating in 60")
// 	time.Sleep(60 * time.Second)

// }

func main() {
	cntxt := &daemon.Context{
		PidFileName: "example.pid",
		PidFilePerm: 0644,
		LogFileName: "ssh-task-runner.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[ssh-task-runner]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("Daemon started")

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker()
	}

	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/add", handleAddTask)
	http.HandleFunc("/stop", handleStop)

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	server := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err.Error())
		}
	}()

	// Listen for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	// Gracefully shut down the server, waiting for existing connections to close

	if err := server.Shutdown(nil); err != nil {
		log.Fatalf("Could not gracefully shut down the server: %s\n", err)
	}

	// Wait for all workers to finish
	wg.Wait()

	// Save results before exiting
	saveResults()
}
