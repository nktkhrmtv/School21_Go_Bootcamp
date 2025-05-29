package main

import (
	"container/heap"
	"errors"
	"fmt"
)

type Present struct {
	Value int
	Size  int
}

type PresentHeap []Present

func (h PresentHeap) Len() int {
	return len(h)
}

func (h PresentHeap) Less(i, j int) bool {
	if h[i].Value == h[j].Value {
		return h[i].Size < h[j].Size
	}
	return h[i].Value > h[j].Value
}

func (h PresentHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *PresentHeap) Push(x interface{}) {
	*h = append(*h, x.(Present))
}

func (h *PresentHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}


func getNCoolestPresents(presents []Present, n int) ([]Present, error) {
	if n > len(presents) || n < 0 {
		return nil, errors.New("n is invalid")
	}
	h := &PresentHeap{}
	heap.Init(h)

	for _, p := range presents {
		heap.Push(h, p)
	}

	result := make([]Present, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, heap.Pop(h).(Present))
	}

	return result, nil
}

func main() {
	presents := []Present{
		{Value: 5, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2},
	}
	presents = append(presents, Present{Value: 15, Size: 1})
	presents = append(presents, Present{Value: 1, Size: 1})

	coolestPresents, err := getNCoolestPresents(presents, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Самые крутые подарки:", coolestPresents)
}