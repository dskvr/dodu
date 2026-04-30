package docker

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinel errors returned by Client implementations.
var (
	// ErrDaemonUnreachable indicates the Docker daemon could not be contacted.
	ErrDaemonUnreachable = errors.New("docker daemon unreachable")
	// ErrPermissionDenied indicates the current user cannot access the daemon socket.
	ErrPermissionDenied = errors.New("permission denied accessing docker socket")
	// ErrReadOnly indicates a destructive call was attempted while DODU_READONLY=1.
	ErrReadOnly = errors.New("read-only mode: destructive operations disabled")
)

// mapError translates a raw SDK / network error into one of the sentinel errors
// where possible, keeping the original message for context.
func mapError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	low := strings.ToLower(msg)
	switch {
	case strings.Contains(low, "permission denied"):
		return fmt.Errorf("%w: %s", ErrPermissionDenied, msg)
	case strings.Contains(low, "no such file or directory") && strings.Contains(low, "docker.sock"):
		return fmt.Errorf("%w: %s", ErrDaemonUnreachable, msg)
	case strings.Contains(low, "connection refused"),
		strings.Contains(low, "cannot connect"),
		strings.Contains(low, "no such host"),
		strings.Contains(low, "is the docker daemon running"):
		return fmt.Errorf("%w: %s", ErrDaemonUnreachable, msg)
	default:
		return err
	}
}
