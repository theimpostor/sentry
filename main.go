package main

import (
	"log"
	// "time"

	"os"
	"os/exec"

	"github.com/rjeczalik/notify"
)

func runCmd(name string, args []string) {
	cmd := exec.Command(name, args...)

	// Attach stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

func main() {
	// Make sure we have arguments, otherwise exit
	if len(os.Args) < 2 {
		panic("Please provide a command to execute")
	}

	// Get the command and arguments from os.Args
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening for events within a directory tree rooted
	// at current working directory. Dispatch remove events to c.
	if err := notify.Watch("./...", c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	// initial run
	runCmd(cmdName, cmdArgs)

	for {
		// Block until an event is received.
		ei := <-c
		log.Println("Got event:", ei)

		runCmd(cmdName, cmdArgs)
	}
}
