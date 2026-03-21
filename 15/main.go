package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func cmdPs() ([]string, error) {
	var out []string

	out = append(out, fmt.Sprintf("%-8s %-20s", "PID", "PROGRAM"))

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return []string{}, fmt.Errorf("failed to read /proc: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		comm, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "comm"))
		if err != nil {
			continue
		}

		program := strings.TrimSpace(string(comm))

		out = append(out, fmt.Sprintf("%-8d %-20s", pid, program))
	}

	return out, nil
}

func cmdCd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("cd: missing argument")
	}
	return os.Chdir(args[1])
}

func cmdKill(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("kill: missing arguments")
	}
	pid, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("kill: invalid PID")
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("kill: cant find process")
	}

	if err := process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}
	return nil
}

func ParseCommand(args []string, r io.Reader, w io.Writer) error {
	comm := strings.ToLower(args[0])
	switch comm {
	case "cd":
		return cmdCd(args)
	case "pwd":
		dir, err := os.Getwd()
		if errors.Is(err, nil) {
			fmt.Fprintln(w, dir)
		}
		return err
	case "echo":
		fmt.Fprintln(w, strings.Join(args[1:], " "))
	case "kill":
		if err := cmdKill(args); err != nil {
			return err
		}
	case "ps":
		pids, err := cmdPs()
		if err != nil {
			return fmt.Errorf("ps error, %w", err)
		}
		for _, pid := range pids {
			fmt.Fprintln(w, pid)
		}
		return nil
	case "exit":
		fmt.Fprintln(w, "Bye Bye!")
		os.Exit(0)
	default:
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = r
		cmd.Stdout = w
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
		w, r := os.Stdout, os.Stdin
		fmt.Print("SfilyaShell>> ")

		if !scanner.Scan() {
			fmt.Println("\nBye")
			break
		}
		line := scanner.Text()
		if line == "" {
			continue
		}

		line = strings.TrimSpace(line)
		input := strings.Split(line, "|")
		if len(input) == 1 {
			if err := ParseCommand(strings.Fields(input[0]), r, w); err != nil {
				slog.Error("err parse command", "error", err)
			}
		}

	}
}
