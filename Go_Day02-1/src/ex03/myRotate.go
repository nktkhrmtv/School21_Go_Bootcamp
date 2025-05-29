package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	archiveDir := flag.String("a", "", "Directory to store archived files")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Error: No files specified.")
		return
	}
	if *archiveDir == "" {
		*archiveDir = "."
	}

	if err := os.MkdirAll(*archiveDir, 0755); err != nil {
		fmt.Println("Error creating archive directory:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		go func(file string) {
			defer wg.Done()

			info, err := os.Stat(file)
			if err != nil {
				fmt.Printf("Error getting file info for %s: %v\n", file, err)
				return
			}

			timestamp := info.ModTime().Unix()
			baseName := filepath.Base(file)
			archiveName := fmt.Sprintf("%s_%d.tar.gz", baseName, timestamp)
			archivePath := filepath.Join(*archiveDir, archiveName)

			if err := createArchive(*archiveDir, archiveName, file); err != nil {
				fmt.Printf("Error creating archive for %s: %v\n", file, err)
				return
			}

			fmt.Printf("Archived %s to %s\n", file, archivePath)
		}(file)
	}
	wg.Wait()
}

func createArchive(archiveDir, fileArchiveName, file string) error {
	cmd := exec.Command("tar", "-czf", fileArchiveName, "-C", ".", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create archive: %v", err)
	}

	exec.Command("mv", fileArchiveName, archiveDir).Run()
	return nil
}
