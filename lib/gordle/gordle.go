package gordle

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

var ErrWordNotFound = errors.New("word not found with requested length")

const (
	FeedbackNotInWord  = Feedback("N")
	FeedbackWrongPlace = Feedback("W")
	FeedbackCorrect    = Feedback("C")
)

type Feedback string

func GetWord(reader *strings.Reader, wordLength int) (string, error) {
	var candidates []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) == wordLength {
			candidates = append(candidates, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("reading from scanner: %w", err)
	}

	candidatesCount := len(candidates)
	if candidatesCount == 0 {
		return "", ErrWordNotFound
	}

	rnd := rand.Intn(candidatesCount) //nolint:gosec // No need for crypto here.

	return candidates[rnd], nil
}

func CheckWord(word, input string) []Feedback {
	wordRunes := []rune(word)
	inputRunes := []rune(input)
	res := []Feedback{}

	counter := map[rune]int{}

	for _, r := range wordRunes {
		counter[r]++
	}

	for i, wordRune := range wordRunes {
		if i >= len(inputRunes) {
			res = append(res, FeedbackNotInWord)
			continue
		}

		inputRune := inputRunes[i]

		switch {
		case wordRune == inputRune:
			res = append(res, FeedbackCorrect)
			counter[wordRune]--
		case counter[inputRune] > 0:
			res = append(res, FeedbackWrongPlace)
			counter[inputRune]--
		default:
			res = append(res, FeedbackNotInWord)
		}
	}

	return res
}
