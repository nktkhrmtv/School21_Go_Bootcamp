package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	linesFlag, charsFlag, wordsFlag := parseFlags()
	files := flag.Args()

	if len(files) == 0 {
		fmt.Println("Error: No files specified.")
		return
	}

	results := processFiles(files, linesFlag, charsFlag, wordsFlag)
	for res := range results {
		fmt.Println(res)
	}
}

func parseFlags() (bool, bool, bool) {
	linesFlag := flag.Bool("l", false, "Count lines")
	charsFlag := flag.Bool("m", false, "Count characters")
	wordsFlag := flag.Bool("w", false, "Count words")
	flag.Parse()

	if (*linesFlag && *charsFlag) || (*linesFlag && *wordsFlag) || (*charsFlag && *wordsFlag) {
		fmt.Println("Error: Only one of -l, -m, or -w can be specified.")
		os.Exit(1)
	}

	if !*linesFlag && !*charsFlag && !*wordsFlag {
		*wordsFlag = true
	}

	return *linesFlag, *charsFlag, *wordsFlag
}

func processFiles(files []string, linesFlag, charsFlag, wordsFlag bool) <-chan string {
	results := make(chan string, len(files))
	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {
		go func(file string) {
			defer wg.Done()
			lines, words, chars := countStats(file)
			result := formatResult(lines, words, chars, linesFlag, charsFlag, wordsFlag, file)
			results <- result
		}(file)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func countStats(file string) (int, int, int) {
	f, err := os.Open(file)
	if err != nil {
		return 0, 0, 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines, words, chars int
	for scanner.Scan() {
		line := scanner.Text()
		lines++
		words += len(strings.Fields(line))
		chars += len([]rune(line))
	}

	return lines, words, chars
}

func formatResult(lines, words, chars int, linesFlag, charsFlag, wordsFlag bool, file string) string {
	var result int
	switch {
	case linesFlag:
		result = lines
	case charsFlag:
		result = chars
	case wordsFlag:
		result = words
	}

	return fmt.Sprintf("%d\t%s", result, file)
}