// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

const (
	Running  = string('\U0001F3C3')
	Complete = string('\U0001F600')
	Failed   = string('\U0001F4A5')
)

// Set a PRINT_FAILURES env variable to 1 to enable
var printFailLogs bool

type taskStatus struct {
	status  string
	logFile string
	name    string
}

func printTaskStatus(tasks []taskStatus, update bool) {
	if update {
		// Clear the last update
		for i := 0; i < len(tasks); i++ {
			fmt.Print("\033[1A") // Up a line
			fmt.Print("\033[K")  // Clear it
		}
	}

	for _, task := range tasks {
		outputLocation := ""
		if task.status == Failed {
			outputLocation = fmt.Sprintf("- %s", task.logFile)
		}
		fmt.Printf("%s - %s %s\n", task.status, task.name, outputLocation)
	}
}

func main() {
	targets := os.Args[1:]
	if len(targets) == 0 {
		fmt.Println("Please pass list of targets to run")
		return
	}

	if os.Getenv("PRINT_FAILURES") == "1" {
		printFailLogs = true
	}

	tasks := make([]taskStatus, len(targets))
	timeStamp := time.Now().Format("20060102")
	wg := sync.WaitGroup{}
	guard := make(chan struct{}, runtime.NumCPU())
	complete := make(chan bool)

	fmt.Printf("Running targets: %v\n", targets)
	for i, target := range targets {
		wg.Add(1)
		guard <- struct{}{}

		go func(tasks []taskStatus, target string, i int) {
			tasks[i].name = target
			tasks[i].status = Running
			defer wg.Done()

			f, _ := os.CreateTemp("", fmt.Sprintf("%s-%s", target, timeStamp))
			tasks[i].logFile = f.Name()

			cmd := exec.Command("make", target)
			out, err := cmd.CombinedOutput()
			if err != nil {
				tasks[i].status = Failed

				// Only capture the output for failures
				_, _ = f.Write(out)
				if err != nil {
					_, _ = f.WriteString(err.Error())
				}
			} else {
				tasks[i].status = Complete
				defer os.Remove(f.Name())
			}

			<-guard
		}(tasks, target, i)
	}

	time.Sleep(500 * time.Millisecond)
	printTaskStatus(tasks, false)

	go func() {
		// Loop through and print out status until complete
		for {
			select {
			case <-complete:
				printTaskStatus(tasks, true)
				return
			default:
				printTaskStatus(tasks, true)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	wg.Wait()
	complete <- true

	fmt.Println("\nAll tasks complete!")

	hasFailures := false
	for _, task := range tasks {
		if task.status == Failed {
			hasFailures = true

			if printFailLogs {
				contents, _ := os.ReadFile(task.logFile)
				fmt.Printf("\n%s failure output:\n", task.name)
				fmt.Print(string(contents))
				fmt.Printf("\n\n")
			}
		}
	}

	if hasFailures {
		os.Exit(1)
	}
}
