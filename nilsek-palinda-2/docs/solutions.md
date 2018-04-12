# Assignment 2

### Task 1
Explain what is wrong in the code below, and then fix the code so that all data really passes through the channel and gets printed.

#### Bug 1
```Go
package main

import "fmt"

// I want this program to print "Hello world!", but it doesn't work.
func main() {
    ch := make(chan string)
    ch <- "Hello world!"
    fmt.Println(<-ch)
}
```
See: [bug01.go](code/bug01.go) for source code to modify.

#### Solution
Channels are used by _goroutines_ to pass data. That means that in order to send and receive values with channels it has to be utilized by a _goroutine_.
The code for above programme has no goroutine, and therefor an error is thrown: ```fatal error: all goroutines are asleep - deadlock!``` 

To solve this we simply introduce a goroutine to the code: 

```Go
package main

import "fmt"

// I want this program to print "Hello world!", but it doesn't work.
func main() {
	ch := make(chan string)
	go func() {
		ch <- "Hello world!"
	}()
	fmt.Println(<-ch)
}
```
output: 
```
Hello world!
```
The new code sends "Hello world!" from a _goroutine_ through our created channel and "Hello world!" is now printed correctly.

#### Bug 2
```Go
package main

import "fmt"

// This program should go to 11, but sometimes it only prints 1 to 10.
func main() {
    ch := make(chan int)
    go Print(ch)
    for i := 1; i <= 11; i++ {
        ch <- i
    }
    close(ch)
}

// Print prints all numbers sent on the channel.
// The function returns when the channel is closed.
func Print(ch <-chan int) {
    for n := range ch { // reads from channel until it's closed
        fmt.Println(n)
    }
}
```

Here we get a data race, where if `Print()` is too slow, there's a risk of it missing the last int (11). This will happen if the channel in `main()` is closed before the value has been printed by `Print()`. To solve this we make use of the sync package's WaitGroup function.
WaitGroup allows us to sync the _goroutines_ created, making sure all of them are executed. 

Our new code will look like this:

```Go
package main

import (
	"fmt"
	"time"
	"sync"
	)

var wg sync.WaitGroup // initialize the WaitGroup

// This program should go to 11, but sometimes it only prints 1 to 10.
func main() {
	ch := make(chan int)
	go Print(ch)
	for i := 1; i <= 11; i++ {
		wg.Add(1) // Add to the waitgroup, making sure the print is 
		ch <- i
	}
	wg.Wait() // Wait for all the routines to be marked as "Done"
	close(ch) // Close the channel
}

// Print prints all numbers sent on the channel.
// The function returns when the channel is closed.
func Print(ch <-chan int) {
	for n := range ch { // reads from channel until it's closed
		time.Sleep(1000) // For easier noticing the bug
		fmt.Println(n)
		wg.Done() // Mark routine as "Done"
	}
}
```  

What we do is that for `i` in our for loop we add to our `WaitGroup`, telling our program that we need to wait a new routine to finish.
When we have printed the value in our `Print()` function we can mark it as done. All in all `WaitGroup` will know that we need to wait for 11 procedures to end before we stop waiting, allowing us to print all the 11 ints. 

### Task 2
* What happens if you switch the order of the statements `wgp.Wait()` and `close(ch)` in the end of the `main` function?

What will happen is that we are allowed to close the channel before we have waited for all `WaitingGroup`'s to finish. This may lead to that we do not print all of the strings since the channel may close early. I think this may throw an error as we in that case will send data on a closed channel, but not entirely sure.

*What actually happened*: We got an error due to sending on a closed channel.

* What happens if you move the `close(ch)` from the `main` function and instead close the channel in the end of the function `Produce`?

We will then close the channel once a `Produce` finish running. Our other `goroutines` will then try to send on a closed channel. Until the channel is closed, the program will run and print as it should.

*What actually happened*: Hypothesis correct 

* What happens if you remove the statement `close(ch)` completely?

The program will keep the channel open. It will print all the strings and finish, however the channel will stay open till end of program. It will not cause any problem in this programme but it may do this in more complex and bigger programmes.

*What actually happened*: The programme ran as expected, printing all the outputs

* What happens if you increase the number of consumers from 2 to 4?

I think that the programme will run faster in average with more consumers, since there's a sleep between prints. More `goroutines` printing means that we faster will print all the outputs sent through the channel.

*What actually happened*: After some testruns it looks like it runs faster with two more consumers by aprox 200-400 ms

* Can you be sure that all strings are printed before the program stops?

Ideally I think that there should be a `WaitGroup` for the consumers as well, otherwise the program can stop while a `Consume` `goroutine` is still in progress.

Modified code with as follows:

```Go
// Stefan Nilsson 2013-03-13

// This is a testbed to help you understand channels better.
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Use different random numbers each time this program is executed.
	rand.Seed(time.Now().Unix())

	const strings = 32
	const producers = 4
	const consumers = 2

	before := time.Now()
	ch := make(chan string)
	wgp := new(sync.WaitGroup)
	wgc := new(sync.WaitGroup)
	wgp.Add(producers)
	wgc.Add(strings) // All the strings printed by consumers.
	for i := 0; i < producers; i++ {
		go Produce("p"+strconv.Itoa(i), strings/producers, ch, wgp)
	}
	for i := 0; i < consumers; i++ {
		go Consume("c"+strconv.Itoa(i), ch, wgc)
	}
	wgp.Wait() // Wait for all producers to finish.
	wgc.Wait() // Wait for all consumer prints to finish.
	close(ch)
	fmt.Println("time:", time.Now().Sub(before))
}

// Produce sends n different strings on the channel and notifies wg when done.
func Produce(id string, n int, ch chan<- string, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		RandomSleep(100) // Simulate time to produce data.
		ch <- id + ":" + strconv.Itoa(i)
	}
	wg.Done()
}

// Consume prints strings received from the channel until the channel is closed.
func Consume(id string, ch <-chan string, wg *sync.WaitGroup) {
	for s := range ch {
		fmt.Println(id, "received", s)
		RandomSleep(100) // Simulate time to consume data.
		wg.Done()
	}
	
}

// RandomSleep waits for x ms, where x is a random number, 0 â‰¤ x < n,
// and then returns.
func RandomSleep(n int) {
	time.Sleep(time.Duration(rand.Intn(n)) * time.Millisecond)
}
```
### Task 3

Oracle modified Oracle code:

```Go
// Stefan Nilsson 2013-03-13

// This program implements an ELIZA-like oracle (en.wikipedia.org/wiki/ELIZA).
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	star   = "Pythia"
	venue  = "Delphi"
	prompt = "> "
)

func main() {
	fmt.Printf("Welcome to %s, the oracle at %s.\n", star, venue)
	fmt.Println("Your questions will be answered in due time.")

	oracle := Oracle()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fmt.Printf("%s heard: %s\n", star, line)
		oracle <- line // The channel doesn't block.
	}
}

// Oracle returns a channel on which you can send your questions to the oracle.
// You may send as many questions as you like on this channel, it never blocks.
// The answers arrive on stdout, but only when the oracle so decides.
// The oracle also prints sporadic prophecies to stdout even without being asked.
func Oracle() chan<- string {
	questions := make(chan string)
	answers := make(chan string)
	// TODO: Answer questions.
	// TODO: Make prophecies.
	// TODO: Print answers.
	// go getPrediction(questions, answers)
	go prophecy("", answers)
	go makeProphecies(questions, answers)
	go printAnswers(answers)
	return questions
}

func makeProphecies(questions chan string, answers chan string) {
	for {
		select {
		case q := <-questions:
			prophecy(q, answers)
		default:
			prophecy("", answers)
		}
	}
}

func printAnswers(answers chan string) {
	for {
		select {
		case ans := <-answers:
			for _, c := range ans {
				time.Sleep(time.Duration(50) * time.Millisecond)
				fmt.Print(string(c))
			}
		}
		fmt.Print("\n")
	}
}

// This is the oracle's secret algorithm.
// It waits for a while and then sends a message on the answer channel.
// TODO: make it better.
func prophecy(question string, answer chan<- string) {
	// Keep them waiting. Pythia, the original oracle at Delphi,
	// only gave prophecies on the seventh day of each month.
	time.Sleep(time.Duration(20+rand.Intn(10)) * time.Second)

	// Find the longest word.
	longestWord := ""
	words := strings.Fields(question) // Fields extracts the words into a slice.
	for _, w := range words {
		if len(w) > len(longestWord) {
			longestWord = w
		}
	}

	// Answer question or fortell a prophecy, in a given priority order
	insultedAns, _ := regexp.MatchString("(?i)(fuck|damn|pussy)", question)
	sassyAns, _ := regexp.MatchString("(?i)(could you|can you)", question)
	giveAns, _ := regexp.MatchString("(?i)(what|could|answer)", question)
	funnyAns, _ := regexp.MatchString("(?i)(ting goes skra|lava toes|lavatoes)", question)

	if insultedAns {
		answer <- "You are a jittery little thing, are you not."
	} else if sassyAns {
		answer <- "I do not know, can I?"
	} else if giveAns {
		answer <- "The answer you are looking for lies within you."
	} else if funnyAns {
		answer <- "You have your moments. Not many of them, but you have them."
	} else {
		// Can you find the source of the prophecies without cheating? :)
		prediction := []string{
			"Concentrate more on your achievements than your failures.",
			"Once you feel nice about yourself, you have planted the first seed to develop self-confidence.",
			"Your focus determines your reality.",
			"Many of the truths we cling to depend greatly on our own point of view.",
			"I feel his presence. But he can also feel mine. He has come for me.",
			"It is your choice, but I warn you not to underestimate my powers.",
			"Patience, my friend. In time, he will seek *you* out, and when he does, you must bring him before me",
			"Everything is proceeding as I have foreseen.",
			"Train yourself to let go of everything you fear to lose.",
			"Always pass on what you have learned.",
			"Greed can be a very powerful ally.",
			"Sometimes we must let go of our pride and do what is requested of us.",
			"Now, be brave and do not look back. Do not look back.",
			"Who is more foolish? The fool or the fool who follows him?",
			"Your eyes can deceive you; don't trust them.",
			"Remember, concentrate on the moment.",
			"In a dark place we find ourselves, and a little more knowledge lights our way.",
		}
		answer <- longestWord + "... " + prediction[rand.Intn(len(prediction))]
	}
}

func init() { // Functions called "init" are executed before the main function.
	// Use new pseudo random numbers every time.
	rand.Seed(time.Now().Unix())
}


```