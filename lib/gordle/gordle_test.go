package gordle_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/Crocmagnon/gordle/lib/gordle"
)

func TestPickWord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		length  int
		wantErr error
		want    string
	}{
		{length: 4, wantErr: nil, want: "test"},
		{length: 5, wantErr: nil, want: "tests"},
		{length: 6, wantErr: gordle.ErrWordNotFound, want: ""},
	}
	for _, test := range tests {
		test := test
		t.Run(test.want, func(t *testing.T) {
			t.Parallel()
			input := "test\ntests"
			got, err := gordle.PickWord(strings.NewReader(input), test.length)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %q want %q", test.wantErr, err)
			}

			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestCheckWord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		word  string
		input string
		want  gordle.FullFeedback
	}{
		{"empty", "", "", gordle.FullFeedback{}},
		{"missing input", "i", "", gordle.FullFeedback{gordle.FeedbackNotInWord}},
		{"correct single", "i", "i", gordle.FullFeedback{gordle.FeedbackCorrect}},
		{"incorrect place", "ab", "ba", gordle.FullFeedback{gordle.FeedbackWrongPlace, gordle.FeedbackWrongPlace}},
		{
			"some correct some incorrect", "aba", "baa",
			gordle.FullFeedback{gordle.FeedbackWrongPlace, gordle.FeedbackWrongPlace, gordle.FeedbackCorrect},
		},
		{
			"complex", "testing", "xsesing",
			gordle.FullFeedback{
				gordle.FeedbackNotInWord,
				gordle.FeedbackWrongPlace,
				gordle.FeedbackWrongPlace,
				gordle.FeedbackNotInWord,
				gordle.FeedbackCorrect,
				gordle.FeedbackCorrect,
				gordle.FeedbackCorrect,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := gordle.CheckWord(test.word, test.input)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
