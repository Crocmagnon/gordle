package gordle

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

var ErrWordNotFound = errors.New("word not found with requested length")

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
