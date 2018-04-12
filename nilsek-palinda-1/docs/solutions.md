# Assignment 1

<!-- Put your answer here -->

## Exercise 2

#### Loops
>Change the loop condition to stop once the value has stopped changing (or only changes by a very small amount). See if that's more or fewer than 10 iterations. Try other initial guesses for z, like x, or x/2. How close are your function's results to the math.Sqrt in the standard library? 

> _Solution_

```Go
package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0

	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
		fmt.Println(z)
	}
	return z
}

func main() {
	fmt.Println("Own calculation:", Sqrt(2))
	fmt.Println("math.Sqrt:", math.Sqrt(2))
	// fmt.Println("Difference:", Sqrt(2)-math.Sqrt(2))
}
```
> The created function loop differs with |2.220446049250313e-16| from the standard library Sqrt function

#### Slices
>Implement Pic. It should return a slice of length dy, each element of which is a slice of dx 8-bit unsigned integers. When you run the program, it will display your picture, interpreting the integers as grayscale (well, bluescale) values. 

> _Solution_

```Go
package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	pic := make([][]uint8, dy)
	for i := 0; i < dy; i++ {
		pic[i] = make([]uint8, dx)
		for j := 0; j < dx; j++ {
			pic[i][j] = uint8((i ^ j) * ((j + i) / 2))
		}
	}
	return pic
}

func main() {
	pic.Show(Pic)
}
```

#### Maps
>Implement WordCount. It should return a map of the counts of each “word” in the string s. The wc.Test function runs a test suite against the provided function and prints success or failure. 

> _Solution_

```Go
package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	words := make(map[string]int)
	wordArray := strings.Fields(s)
	for i := 0; i < len(wordArray); i++ {
		words[wordArray[i]] += 1
	}
	return words
}

func main() {
	wc.Test(WordCount)
}
```



#### Fibonnaci
> Implement a fibonacci function that returns a function (a closure) that returns successive fibonacci numbers (0, 1, 1, 2, 3, 5, ...). 

> _Solution_

```Go
package main

import "fmt"

func fibonacci() func() int {
	nextf := 0
	f1 := 1
	f2 := 0

	return func() int {
		nextf = f1 + f2
		f1, f2 = f2, nexf
		return f1
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
```

## Exercice 3

>The output will repeatedly print the output after the given delay, and `XX.XX` should be replaced with the current time, and `<text>` should be replaced by the contents of `text`.

> _Solution_

```Go
package remind

import (
	"fmt"
	"time"
)

func Remind(text string, delay time.Duration) {
	for {
		time.Sleep(delay)
		fmt.Println("Klockan är: ", time.Now().Format("15.04"), text)
	}
}

func main() {
	// set duration of delay
	eat := 3 * time.Hour
	work := 8 * time.Hour
	sleep := 24 * time.Hour

	// start threads
	go Remind("Go eat something", eat)
	go Remind("Time to work", work)
	go Remind("Nap time!", sleep)
	select {}
}

```



## Exercice 4

>In this task you will complete the following partial program.  It adds all of the numbers in an array by splitting the array in half, then having two Go routines take care of each half.

> _Solution_

```Go
package main

import "fmt"

// Add adds the numbers in a and sends the result on res.
func Add(a []int, res chan<- int) {
	// TODO: return the sum of the numbers on the given channel
	sum := 0
	for _, i := range a {
		sum += i
	}
	res <- sum
}

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7}
	n := len(a)
	ch := make(chan int)
	go Add(a[:n/2], ch)
	go Add(a[n/2:], ch)

	// TODO: Get the subtotals from the channel and print their sum.
	x, y := <-ch, <-ch
	fmt.Println("Sum of subtotals: ", x+y)
}

```
>_Unittest_

```Code
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

```