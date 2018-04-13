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
			// Add a delay between prophecies, the common people doesn't
			// deserve THAT many prophecies...
			time.Sleep(time.Duration(20+rand.Intn(10)) * time.Second)
			prophecy("", answers)
		}
	}
}

func printAnswers(answers chan string) {
		for ans := range answers {
			for _, c := range ans {
				time.Sleep(time.Duration(50) * time.Millisecond)
				fmt.Print(string(c))
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
			"Many of the tru 	ths we cling to depend greatly on our own point of view.",
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
