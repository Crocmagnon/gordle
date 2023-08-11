package gordle

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

var (
	ErrWordNotFound = errors.New("word not found with requested length")
	ErrGameWon      = errors.New("won game")
	ErrGameLost     = errors.New("game over")
)

type (
	Feedback     string
	FullFeedback []Feedback
)

func (ff FullFeedback) Wins() bool {
	for _, f := range ff {
		if f != FeedbackCorrect {
			return false
		}
	}

	return true
}

func (ff FullFeedback) String() string {
	str := strings.Builder{}
	for _, f := range ff {
		str.WriteString(string(f))
	}

	return str.String()
}

const (
	FeedbackNotInWord  = Feedback("N")
	FeedbackWrongPlace = Feedback("W")
	FeedbackCorrect    = Feedback("C")
)

type Game struct {
	maxAttempts    int
	currentAttempt int
	pickedWord     string
	over           bool
}

func New(maxAttempts int, pickedWord string) *Game {
	g := &Game{
		maxAttempts: maxAttempts,
		pickedWord:  pickedWord,
	}

	return g
}

func (g *Game) TryWord(word string) (FullFeedback, error) {
	feedback := CheckWord(g.pickedWord, word)
	if feedback.Wins() {
		g.over = true
		return feedback, ErrGameWon
	}

	g.currentAttempt++
	if g.currentAttempt >= g.maxAttempts {
		g.over = true
		return feedback, ErrGameLost
	}

	return feedback, nil
}

func (g *Game) Over() bool {
	return g.over
}

func PickWord(reader io.Reader, wordLength int) (string, error) {
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

func CheckWord(word, input string) FullFeedback {
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
