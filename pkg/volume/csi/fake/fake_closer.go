package fake

import (
	"testing"
)

func NewCloser(t *testing.T) *Closer {
	return &Closer{
		t: t,
	}
}

type Closer struct {
	wasCalled bool
	t         *testing.T
}

func (c *Closer) Close() error {
	c.wasCalled = true
	return nil
}

func (c *Closer) Check() *Closer {
	c.t.Helper()

	if !c.wasCalled {
		c.t.Error("expected closer to have been called")
	}

	return c
}
