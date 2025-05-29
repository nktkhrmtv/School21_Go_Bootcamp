package main

import (
    "testing"
    "day07/minCoinsFuncs"
)

type Tests struct{
    val   int
    coins []int
    want  int 
}

func TestMinCoins(t *testing.T) {
    tests := []Tests{
        {6, []int{1, 3, 4}, 2}, 
        {13, []int{1, 5, 10}, 4}, 
        {14, []int{1, 7, 10}, 2},
        {6, []int{1, 2, 5}, 2}, 
        {9, []int{1, 3, 4}, 3}, 
        {0, []int{1, 2, 3}, 0}, 
        {5, []int{}, 0}, 
        {5, []int{1, 1, 2, 3}, 2}, 
        {90, []int{1, 6, 84, 85}, 2}, 
    }

    for _, tt := range tests {
        got := minCoinsFuncs.MinCoins(tt.val, tt.coins)
        if len(got) != tt.want {
            t.Errorf("minCoins(%d, %v) = %v (количество монет: %d), минимальное количество монет: %d", tt.val, tt.coins, got, len(got), tt.want)
        }

        sum := 0
        for _, coin := range got {
            sum += coin
        }
        if sum != tt.val {
            t.Errorf("minCoins(%d, %v) = %v (сумма: %d), нужная сумма: %d", tt.val, tt.coins, got, sum, tt.val)
        }
    }
}

func TestMinCoins2(t *testing.T) {
    tests := []Tests{
        {6, []int{1, 3, 4}, 2},
        {13, []int{1, 5, 10}, 4}, 
        {14, []int{1, 7, 10}, 2},
        {6, []int{1, 2, 5}, 2}, 
        {9, []int{1, 3, 4}, 3},
        {0, []int{1, 2, 3}, 0}, 
        {5, []int{}, 0}, 
        {5, []int{1, 1, 2, 3}, 2},
        {90, []int{1, 6, 84, 85}, 2}, 
    }

    for _, tt := range tests {
        got := minCoinsFuncs.MinCoins2(tt.val, tt.coins)
        if len(got) > tt.want {
            t.Errorf("minCoins2(%d, %v) = %v (количество монет: %d), минимальное количество монет: %d", tt.val, tt.coins, got, len(got), tt.want)
        }

        sum := 0
        for _, coin := range got {
            sum += coin
        }
        if sum != tt.val && len(tt.coins) != 0{
            t.Errorf("minCoins2(%d, %v) = %v (сумма: %d), нужная сумма: %d", tt.val, tt.coins, got, sum, tt.val)
        }
    }
}