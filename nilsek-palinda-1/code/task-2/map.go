package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	words := make(map[string]int)
	wordArray := strings.Fields(s)
	for _, v := range wordArray {
		words[v]++
	}
	return words
}

func main() {
	wc.Test(WordCount)
}
