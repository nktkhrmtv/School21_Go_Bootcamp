package main

import (
    "os"
    "runtime/pprof"
	"day07/minCoinsFuncs"
)

func main() {
    f, err := os.Create("cpu_profile.prof")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    if err := pprof.StartCPUProfile(f); err != nil {
        panic(err)
    }
    defer pprof.StopCPUProfile()

    testMinCoins2()
}

func testMinCoins2() {
    val := 1000
    coins := []int{1, 3, 4, 7, 13, 15}
    for i := 0; i < 100000; i++ {
        minCoinsFuncs.MinCoins2(val, coins)
    }
}