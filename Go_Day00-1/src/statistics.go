package main

import (
	"math"
	"sort"
)

type Statistics struct {
	numbers []int
	mean   float64
	median float64
	mode   int
	sd     float64
}

func (s *Statistics) Mean() float64 {
	sum := 0
	for _, num := range s.numbers {
		sum += num
	}
	s.mean = float64(sum) / float64(len(s.numbers))
	return s.mean
}

func (s *Statistics) Median() float64 {
	sort.Ints(s.numbers)
	n := len(s.numbers)
	if n % 2 == 0 {
		s.median = float64(s.numbers[n/2-1] + s.numbers[n/2]) / 2.0
	} else {
		s.median = float64(s.numbers[n/2])
	}
	return s.median
}

func (s *Statistics) Mode() int {
	frequency := make(map[int]int)
	for _, num := range s.numbers {
		frequency[num]++
	}
	s.mode = 0
	maxFreq := 0
	for num, freq := range frequency {
		if freq > maxFreq || (freq == maxFreq && num < s.mode) {
			maxFreq = freq
			s.mode = num
		}
	}
	return s.mode
}

func (s *Statistics) SD() float64 {
	m := s.Mean()
	var sum float64
	for _, num := range s.numbers {
		sum += math.Pow(float64(num)-m, 2)
	}
	variance := sum / float64(len(s.numbers))
	s.sd = math.Sqrt(variance)
	return s.sd 
}