package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	oldFile := flag.String("old", "", "Path to the old snapshot file")
	newFile := flag.String("new", "", "Path to the new snapshot file")
	flag.Parse()
	if *oldFile == "" || *newFile == "" {
		fmt.Println("Please provide both --old and --new file paths.")
		return
	}

	oldSnapshot, err := readSnapshot(*oldFile)
	if err != nil {
		fmt.Println("Error reading old snapshot:", err)
		return
	}
	newSnapshot, err := readSnapshot(*newFile)
	if err != nil {
		fmt.Println("Error reading new snapshot:", err)
		return
	}

	compareSnapshots(oldSnapshot, newSnapshot)
}

func readSnapshot(filename string) (map[string]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	snapshot := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := scanner.Text()
		snapshot[path] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return snapshot, nil
}

func compareSnapshots(oldSnapshot, newSnapshot map[string]bool) {
	for path := range newSnapshot {
		if !oldSnapshot[path] {
			fmt.Printf("ADDED %s\n", path)
		}
	}
	for path := range oldSnapshot {
		if !newSnapshot[path] {
			fmt.Printf("REMOVED %s\n", path)
		}
	}
}
