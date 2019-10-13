package lib

import (
	"golang.org/x/xerrors"
)

type StackTracer interface {
	Wrap(string, error) error
}

func NewStackTracer() StackTracer {
	return &stackTracer{}
}

type stackTracer struct{}

func (st *stackTracer) Wrap(message string, err error) error {
	return xerrors.Errorf(message+": %w", err)
}
