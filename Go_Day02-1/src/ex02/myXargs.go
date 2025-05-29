package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <command1> [args1...] ./myXargs <command2> [args2...]")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	var inputArgs []string
	for scanner.Scan() {
		inputArgs = append(inputArgs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
		return
	}

	command := os.Args[1]
	args := os.Args[2:]
	args = append(args, inputArgs...)

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}