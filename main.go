package main

import (
	"log"
	// "time"

	"flag"
	"os"
	"os/exec"

	"github.com/rjeczalik/notify"
)

func runCmd(name string, args []string) {
	cmd := exec.Command(name, args...)

	// Attach stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Running", name)

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Println(name, "exited with error", exitError.ExitCode())
		} else {
			log.Fatal(err)
		}
	} else {
		log.Println(name, "finished successfully")
	}
}

func main() {
	// Add a flag to specify the watched directory, otherwise use the current directory
	var watchDir string
	flag.StringVar(&watchDir, "d", ".", "Directory to watch")
	flag.Parse()

	// Make sure we have arguments, otherwise exit
	if len(flag.Args()) < 1 {
		log.Fatal("Please provide a command to execute")
	}

	log.Printf("Command to run: %+q\n", flag.Args())

	cmdName := flag.Args()[0]
	cmdArgs := flag.Args()[1:]

	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening for events within a directory tree rooted
	// at watchDir. Dispatch remove events to c.
	if err := notify.Watch(watchDir+"/...", c, notify.All); err != nil {
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
