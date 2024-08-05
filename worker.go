package main

import (
	"log"
	"time"
)

// var (
// 	taskQueue   = make(chan Task, 100)
// 	results     = sync.Map{}
// 	taskCounter = 0
// 	totalTasks  = 0
// 	wg          sync.WaitGroup
// 	stop        = make(chan struct{})
// 	ssh         = alfredo.SSHStruct{Key: alfredo.ExpandTilde("~/.ssh/homelab_rsa"), Host: "192.168.1.10", User: "cdelezenski"}
// )

// func initTasks(numTasks int) {
// 	var r int
// 	for i := 0; i < numTasks; i++ {
// 		r = rand.Intn(20) + 10
// 		taskQueue <- Task{ID: i, ssh: ssh, IsRunning: false, Command: fmt.Sprintf("sleep %d; echo Task %d", r, i)}
// 	}
// 	close(taskQueue)
// }

// func startWorkers(numWorkers int) {
// 	for i := 0; i < numWorkers; i++ {
// 		wg.Add(1)
// 		go worker()
// 	}
// }

// func waitForCompletion() {
// 	wg.Wait()
// }

// func worker() {
// 	defer wg.Done()

// 	for task := range taskQueue {
// 		task.StartTime = time.Now()

// 		log.Printf("about to start task %d, storing result in result map", task.ID)
// 		task.IsRunning = true
// 		results.Store(task.ID, task)

// 		err := task.ssh.SecureRemoteExecution(task.Command)

// 		task.Duration = time.Since(task.StartTime)
// 		task.Output = task.ssh.GetBody()
// 		if err != nil {
// 			task.Error = err.Error()
// 		} else {
// 			task.Error = ""
// 		}
// 		task.IsRunning = false
// 		results.Store(task.ID, task)
// 		log.Printf("Task stored: %+v\n", task) // Logging task details
// 	}
// }
func worker() {
	defer wg.Done()
	for {
		select {

		case task := <-taskQueue:

			log.Printf("should be updating queuedTasks --1 ")
			// Remove task from queuedTasks
			for i, t := range queuedTasks {
				if t.ID == task.ID {
					log.Printf("should be updating queuedTask --2")
					queuedTasks = append(queuedTasks[:i], queuedTasks[i+1:]...)
					break
				}
			}
			//log.Printf("\t")
			if len(task.Command) == 0 {
				log.Fatalf("worker: blank command requested")
				break
			}
			log.Printf("about to start task %d, cli=%q, storing result in result map", task.ID, task.Command)
			task.IsRunning = true
			task.StartTime = time.Now()

			results.Store(task.ID, task)

			if len(task.ssh.Key) == 0 {
				panic("ssh key is blank, shouldn't be")
			}

			err := task.ssh.SecureRemoteExecution(task.Command)

			task.Duration = time.Since(task.StartTime)
			task.Output = task.ssh.GetBody()
			if err != nil {
				task.Error = err.Error()
			} else {
				task.Error = ""
			}
			task.IsRunning = false
			task.HasRun = true
			results.Store(task.ID, task)
			log.Printf("Task stored: %+v\n", task) // Logging task details

		case <-stop:
			return
		}

	}
}

// // executeSSHCommand executes the given command over SSH and returns the output.
// func executeSSHCommand(command string) (string, error) {

// 	err := s.SecureRemoteExecution(command)

// 	// Dummy SSH execution. Replace with actual SSH logic.
// 	return s.GetBody(), err
// }
