package main

import "testing"

type testpair struct {
	values     []int
	twopartsum int
}

var tests = []testpair{
	{[]int{1, 1, 1, 2, 2, 2}, 9},
	{[]int{3, 2, 1, 1, 2, 4, 3}, 16},
	{[]int{3, 2, 2, 7}, 14},
}

var tests2 = []testpair{
	{[]int{1, 1, 1, 2, 2, 2}, 3},
	{[]int{3, 2, 1, 1, 2, 4, 3}, 6},
	{[]int{3, 2, 2, 7}, 5},
}

var tests3 = []testpair{
	{[]int{1, 1, 1, 2, 2, 2}, 6},
	{[]int{3, 2, 1, 1, 2, 4, 3}, 10},
	{[]int{3, 2, 2, 7}, 9},
}

func TestTwopartsum(t *testing.T) {
	ch := make(chan int)
	for _, pair := range tests {
		n := len(pair.values)
		go Add(pair.values[:n/2], ch)
		go Add(pair.values[n/2:], ch)
		var x, y = <-ch, <-ch

		if x+y != pair.twopartsum {
			t.Error(
				"For", pair.values,
				"expected", pair.twopartsum,
				"got", x+y,
			)
		}
	}
}
