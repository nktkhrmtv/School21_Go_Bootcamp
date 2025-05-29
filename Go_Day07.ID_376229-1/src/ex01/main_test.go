package main

import (
    "testing"
	"day07/minCoinsFuncs"
)

func BenchmarkMinCoins(b *testing.B) {
    val := 1000
    coins := []int{1, 3, 4, 7, 13, 15}
    for i := 0; i < b.N; i++ {
        minCoinsFuncs.MinCoins(val, coins)
    }
}

func BenchmarkMinCoins2(b *testing.B) {
    val := 1000
    coins := []int{1, 3, 4, 7, 13, 15}
    for i := 0; i < b.N; i++ {
        minCoinsFuncs.MinCoins2(val, coins)
    }
}