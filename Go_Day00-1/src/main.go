package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func Init(numbers []int) *Statistics {
	return &Statistics{
		numbers: numbers,
	}
}

func main() {
	meanFlag := flag.Bool("mean", false, "Print mean value")
	medianFlag := flag.Bool("median", false, "Print median value")
	modeFlag := flag.Bool("mode", false, "Print mode value")
	sdFlag := flag.Bool("sd", false, "Print standard deviation value")
	flag.Parse()

	var numbers []int
	err := readNumbers(&numbers)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading numbers:", err)
		return
	}

	stats := Init(numbers)

	if *meanFlag || *medianFlag || *modeFlag || *sdFlag {
		if *meanFlag {
			fmt.Printf("Mean: %.2f\n", stats.Mean())
		}
		if *medianFlag {
			fmt.Printf("Median: %.2f\n", stats.Median())
		}
		if *modeFlag {
			fmt.Printf("Mode: %d\n", stats.Mode())
		}
		if *sdFlag {
			fmt.Printf("SD: %.2f\n", stats.SD())
		}
	} else {
		fmt.Printf("Mean: %.2f\n", stats.Mean())
		fmt.Printf("Median: %.2f\n", stats.Median())
		fmt.Printf("Mode: %d\n", stats.Mode())
		fmt.Printf("SD: %.2f\n", stats.SD())
	}
}

func readNumbers(numbers *[]int) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter numbers (one per line). Type 'q' to finish:")
	var err error
	for scanner.Scan() {
		input := scanner.Text()
		if input == "q" {
			break
		}
		var num int
		num, err = strconv.Atoi(input)
		if err != nil {
			err = fmt.Errorf("invalid input: %s is not a number", input)
			break
		}
		if num < -100000 || num > 100000 {
			err = fmt.Errorf("number %d is out of range (-100000 to 100000)", num)
			break
		}
		*numbers = append(*numbers, num)
	}
	if err == nil && len(*numbers) == 0 {
		err = errors.New("no numbers provided")
	}
	return err
}