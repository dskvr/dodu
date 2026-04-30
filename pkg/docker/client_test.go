package docker

import (
	"testing"
	"time"
)

func TestDefaultOptions(t *testing.T) {
	o := defaultOptions()
	if o.timeout == 0 {
		t.Errorf("expected non-zero default timeout")
	}
}

func TestWithHostAndTimeout(t *testing.T) {
	o := defaultOptions()
	WithHost("tcp://example:2375")(&o)
	WithTimeout(5 * time.Second)(&o)
	if o.host != "tcp://example:2375" {
		t.Errorf("host = %q", o.host)
	}
	if o.timeout != 5*time.Second {
		t.Errorf("timeout = %v", o.timeout)
	}
}
