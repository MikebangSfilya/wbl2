package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func fields(line string) []string {
	return strings.Fields(line)

}

func ParseCommand(args []string) error {
	comm := strings.ToLower(args[0])
	switch comm {
	case "cd":
		if len(args) < 2 {
			return fmt.Errorf("cd: missing argument")
		}
		return os.Chdir(args[1])
	case "pwd":
		dir, err := os.Getwd()
		if errors.Is(err, nil) {
			fmt.Println(dir)
		}
		return err
	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
	case "exit":
		os.Exit(0)
	default:
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("SfilyaShell: command not found: %s\n", args[0])
		}
	}
	return nil
}

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		for {
			<-sigs
			fmt.Print("\nInterrupt received (Ctrl+C). Type 'exit' to quit.\nSfilyaShell$ ")
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to Sfilya MiniShell, wb l2.15")
	for {
		fmt.Print("SfilyaShell$ ")

		if !scanner.Scan() {
			fmt.Println("\nBye")
			break
		}
		line := scanner.Text()
		if line == "" {
			continue
		}

		args := fields(line)
		ParseCommand(args)
	}
}
