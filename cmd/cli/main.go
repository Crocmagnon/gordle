//nolint:forbidigo
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Crocmagnon/gordle/lib/gordle"
)

const (
	defaultWordLength  = 5
	defaultMaxAttempts = 5
)

func main() {
	wordLength := flag.Int("l", defaultWordLength, "word length")
	dictionary := flag.String("f", "", "dictionary file")
	maxAttempts := flag.Int("m", defaultMaxAttempts, "maximum number of attempts")
	cheat := flag.Bool("c", false, "cheat mode (prints word at the beginning)")

	flag.Parse()

	if !validateInput(dictionary, wordLength, maxAttempts) {
		os.Exit(1)
	}

	file, err := os.Open(*dictionary)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	word, err := gordle.PickWord(file, *wordLength)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("length %d ; you have %d attempts\n", *wordLength, *maxAttempts)

	if *cheat {
		fmt.Printf("word: %s\n", word)
	}

	scanner := bufio.NewScanner(os.Stdin)
	game := gordle.New(*maxAttempts, word)

	play(game, scanner, word, wordLength)
}

func play(game *gordle.Game, scanner *bufio.Scanner, word string, wordLength *int) {
	for !game.Over() && scanner.Scan() {
		var feedback gordle.FullFeedback

		text := strings.ToUpper(scanner.Text())
		feedback, err := game.TryWord(text)

		switch {
		case errors.Is(err, gordle.ErrGameWon):
			fmt.Println(feedback)
			fmt.Println("🎉 you won")
		case errors.Is(err, gordle.ErrGameLost):
			fmt.Println(feedback)
			fmt.Printf("😔 you lost, the correct word was %s\n", word)
		case errors.Is(err, gordle.ErrRunesCount):
			fmt.Printf("🤔 please provide a %d-letters word\n", *wordLength)
		default:
			fmt.Println(feedback)
		}
	}
}

func validateInput(dictionary *string, wordLength *int, maxAttempts *int) bool {
	ok := true //nolint:varnamelen

	if *dictionary == "" {
		fmt.Println("dictionary file is mandatory")

		ok = false
	}

	if *wordLength <= 0 {
		fmt.Println("word length must be positive")

		ok = false
	}

	if *maxAttempts <= 0 {
		fmt.Println("max attempts must be positive")

		ok = false
	}

	return ok
}
