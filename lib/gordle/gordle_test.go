package gordle_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Crocmagnon/gordle/lib/gordle"
)

func TestGetWord(t *testing.T) {
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
			got, err := gordle.GetWord(strings.NewReader(input), test.length)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("got err %q want %q", test.wantErr, err)
			}

			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
