package docker

import (
	"errors"
	"strings"
	"testing"
)

func TestMapErrorPermissionDenied(t *testing.T) {
	in := errors.New("Got permission denied while trying to connect to the Docker daemon socket")
	got := mapError(in)
	if !errors.Is(got, ErrPermissionDenied) {
		t.Fatalf("expected ErrPermissionDenied, got %v", got)
	}
	if !strings.Contains(got.Error(), "permission denied") {
		t.Errorf("expected original message preserved, got %q", got.Error())
	}
}

func TestMapErrorDaemonUnreachable(t *testing.T) {
	cases := []string{
		"Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?",
		"dial unix /var/run/docker.sock: connect: connection refused",
		"open /var/run/docker.sock: no such file or directory",
	}
	for _, msg := range cases {
		t.Run(msg, func(t *testing.T) {
			got := mapError(errors.New(msg))
			if !errors.Is(got, ErrDaemonUnreachable) {
				t.Errorf("expected ErrDaemonUnreachable for %q, got %v", msg, got)
			}
		})
	}
}

func TestMapErrorPassthrough(t *testing.T) {
	in := errors.New("some unrelated failure")
	got := mapError(in)
	if errors.Is(got, ErrDaemonUnreachable) || errors.Is(got, ErrPermissionDenied) {
		t.Errorf("did not expect sentinel for %v", got)
	}
	if got.Error() != in.Error() {
		t.Errorf("expected passthrough, got %q", got.Error())
	}
}

func TestMapErrorNil(t *testing.T) {
	if got := mapError(nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}
