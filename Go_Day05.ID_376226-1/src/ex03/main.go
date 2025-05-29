package main

import (
	"fmt"
)

type Present struct {
	Value int
	Size  int
}

func grabPresents(presents []Present, capacity int) []Present {
	n := len(presents)

	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for j := 0; j <= capacity; j++ {
			if presents[i-1].Size > j {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-presents[i-1].Size]+presents[i-1].Value)
			}
		}
	}
	result := make([]Present, 0)
	j := capacity
	for i := n; i > 0; i-- {
		if dp[i][j] != dp[i-1][j] {
			result = append(result, presents[i-1])
			j -= presents[i-1].Size
		}
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	presents := []Present{
		{Value: 2, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
		{Value: 1, Size: 3},
		{Value: 1, Size: 1},
		{Value: 8, Size: 1},
	}

	capacity := 3

	result := grabPresents(presents, capacity)
	fmt.Println("Подарки с максимальной стоимостью:", result)
}