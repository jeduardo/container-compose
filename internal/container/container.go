package container

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

// Run a named container
func Run(containerName string, image string) {
	cmd := exec.Command("/usr/local/bin/container", "run", "--name", containerName, image)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("could not get stdout: %s", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("could not get stdout: %s", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("unable to run command: %s", err)
	}

	// stream stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Printf("%s: %s\n", containerName, m)
		}
	}()

	// stream stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Printf("%s: %s\n", containerName, m)
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("exit err:", err)
	}
}
