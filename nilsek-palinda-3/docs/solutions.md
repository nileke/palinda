# Assignment 3

> ### Task 1 - Matching Behaviour

Take a look at the program [matching.go](code/matching.go). Explain what happens and why it happens if you make the following changes. Try first to reason about it, and then test your hypothesis by changing and running the program.

Hint: Think about the order of the instructions and what happens with arrays of different lengths.

<!-- Put your answer here -->
#### Answers

  * What happens if you remove the `go-command` from the `Seek` call in the `main` function?

 > Seek will not be called in a `goroutine` and instead will be run one at a time, program will wait for every Seek-call to finish.

   * What happens if you switch the declaration `wg := new(sync.WaitGroup`) to `var wg sync.WaitGroup` and the parameter `wg *sync.WaitGroup` to `wg sync.WaitGroup`?

 > All the `goroutines` instances of `Seek` can access the same `WaitGroup` safely, marking as done when finished. Changing the code will create a deadlock error, due to the `WaitGroup` not being safely accessed.

   * What happens if you remove the buffer on the channel match?

 > Buffered channels accepts a limited number of values, in this case 1, which means that it doesn't need a corresponding amount of send and receive. When removing the buffering the program would need a corresponding amount of send and receive, which is not the case in our program creating a deadlock. 

  * What happens if you remove the default-case from the case-statement in the `main` function?

 > Default means that if none of the cases criterias are matched the default case is chosen. In this program default is empty so nothing will happen if we remove it.

> ### Task 2 - Fractal Images

The file [julia.go](code/julia.go) contains a program that creates images and writes them to file. The program is pretty slow. Your assignment is to divide the computations so that they run in parallel on all available CPUs. Use the ideas from the example in the [efficient parallel computation](http://yourbasic.org/golang/efficient-parallel-computation/) section of the course literature.

You can also make changes to the program, such as using different functions and other colourings.

How many CPUs does you program use? How much faster is your parallel version?

#### Answers

I modified the code by adding `goroutines` in the mainfunction. By changing the program to run in parallell the computation time went from the average 9.80 seconds to an average of 0.220 seconds. This value did not change noticable by changing the number. 

Modifying Julia to be run with `goroutines` in the for loop did not make any difference in the program's runtime

Modified code:

```Go
// Stefan Nilsson 2013-02-27

// This program creates pictures of Julia sets (en.wikipedia.org/wiki/Julia_set).
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"strconv"
	"sync"
	"fmt"
	"runtime"
)

type ComplexFunc func(complex128) complex128

var Funcs []ComplexFunc = []ComplexFunc{
	func(z complex128) complex128 { return z*z - 0.61803398875 },
	func(z complex128) complex128 { return z*z + complex(0, 1) },
	func(z complex128) complex128 { return z*z + complex(-0.835, -0.2321) },
	func(z complex128) complex128 { return z*z + complex(0.45, 0.1428) },
	func(z complex128) complex128 { return z*z*z + 0.400 },
	func(z complex128) complex128 { return cmplx.Exp(z*z*z) - 0.621 },
	func(z complex128) complex128 { return (z*z+z)/cmplx.Log(z) + complex(0.268, 0.060) },
	func(z complex128) complex128 { return cmplx.Sqrt(cmplx.Sinh(z*z)) + complex(0.065, 0.122) },
}

func main() {
	var wg sync.WaitGroup
	runtime.GOMAXPROCS(4)
	for n, fn := range Funcs {
		go func(n int, fn ComplexFunc) {
			err := CreatePng("picture-"+strconv.Itoa(n)+".png", fn, 1024)
			if err != nil {
				log.Fatal(err)
			}
		wg.Done()
		}(n, fn)
	wg.Wait()
	}
	fmt.Printf("Number of CPUs: " + strconv.Itoa(runtime.NumCPU()))
}

// CreatePng creates a PNG picture file with a Julia image of size n x n.
func CreatePng(filename string, f ComplexFunc, n int) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	err = png.Encode(file, Julia(f, n))
	return
}

// Julia returns an image of size n x n of the Julia set for f.
func Julia(f ComplexFunc, n int) image.Image {
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	s := float64(n / 4)
	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
			r := uint8(0)
			g := uint8(0)
			b := uint8(n % 32 * 8)
			img.Set(i, j, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

// Iterate sets z_0 = z, and repeatedly computes z_n = f(z_{n-1}), n â‰¥ 1,
// until |z_n| > 2  or n = max and returns this n.
func Iterate(f ComplexFunc, z complex128, max int) (n int) {
	for ; n < max; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
		z = f(z)
	}
	return
}
```