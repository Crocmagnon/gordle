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
	ErrRunesCount   = errors.New("incorrect number of letters")
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
	FeedbackNotInWord  = Feedback("⬜️")
	FeedbackWrongPlace = Feedback("🟡")
	FeedbackCorrect    = Feedback("💚")
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
	feedback, err := CheckWord(g.pickedWord, word)
	if err != nil {
		return feedback, err
	}

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

func CheckWord(word, input string) (FullFeedback, error) {
	wordRunes := []rune(word)
	inputRunes := []rune(input)

	if len(wordRunes) != len(inputRunes) {
		return nil, ErrRunesCount
	}

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

	return res, nil
}
