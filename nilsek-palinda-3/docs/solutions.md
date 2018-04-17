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
// Julia returns an image of size n x n of the Julia set for f.
func Julia(f ComplexFunc, n int) image.Image {
	wg := new(sync.WaitGroup)
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	s := float64(n / 4)
	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		wg.Add(1)
		go func(i int) {
			for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
				n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
				r := uint8(0)
				g := uint8(0)
				b := uint8(n % 32 * 8)
				img.Set(i, j, color.RGBA{r, g, b, 255})
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return img
}
```