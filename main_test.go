//go:build !integration

package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"golang.org/x/term"
)

func TestReviewCommitMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		suggest string
		want    string
		ok      bool
		wantErr bool
		errSub  string
	}{
		{
			name:    "yes",
			input:   "y\n",
			suggest: "hello",
			want:    "hello",
			ok:      true,
		},
		{
			name:    "empty line means yes",
			input:   "\n",
			suggest: "hello",
			want:    "hello",
			ok:      true,
		},
		{
			name:    "YES uppercase",
			input:   "YES\n",
			suggest: "hello",
			want:    "hello",
			ok:      true,
		},
		{
			name:    "no aborts",
			input:   "n\n",
			suggest: "hello",
			want:    "",
			ok:      false,
		},
		{
			name:    "invalid then yes",
			input:   "foo\ny\n",
			suggest: "hi",
			want:    "hi",
			ok:      true,
		},
		{
			name:    "eof before newline",
			input:   "",
			suggest: "hi",
			wantErr: true,
			errSub:  "EOF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			in := strings.NewReader(tt.input)
			got, ok, err := reviewCommitMessage(tt.suggest, in, io.Discard)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if tt.errSub != "" && !strings.Contains(err.Error(), tt.errSub) {
					t.Fatalf("error %q should contain %q", err.Error(), tt.errSub)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if ok != tt.ok {
				t.Fatalf("ok: got %v want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("message: got %q want %q", got, tt.want)
			}
		})
	}
}

func TestReviewCommitMessage_editNonTerminal(t *testing.T) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		t.Skip("uses non-terminal edit path; stdin is a TTY")
	}

	in := strings.NewReader("e\nreplaced msg\n")
	got, ok, err := reviewCommitMessage("original", in, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected ok")
	}
	if got != "replaced msg" {
		t.Fatalf("got %q", got)
	}
}
