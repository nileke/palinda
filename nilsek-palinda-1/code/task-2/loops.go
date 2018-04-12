package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	prev := 0.0
	for {
		if math.Abs(z-prev) < 0.0000000001 {
			break
		}
		prev = z
		z -= (z*z - x) / (2*z)
		fmt.Println(z)
	}
	return z
}



func main() {
	fmt.Println("Own calculation:", Sqrt(2))
	fmt.Println("math.Sqrt:", math.Sqrt(2))
	fmt.Println("Difference:", Sqrt(2)-math.Sqrt(2))
}
