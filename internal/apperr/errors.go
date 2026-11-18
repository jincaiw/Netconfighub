// Package apperr defines typed/sentinel errors shared across packages so that
// error classification never has to rely on matching human-readable text.
package apperr

import (
	"errors"
	"fmt"
)

// Sentinel errors. They are compared with errors.Is; their text is intentionally
// not part of any HTTP response or user-facing message.
var (
	ErrConnectTimeout = errors.New("connect timeout")
	ErrAuthFailed     = errors.New("authentication failed")
	ErrPromptTimeout  = errors.New("prompt timeout")
	ErrCommandTimeout = errors.New("command timeout")
	ErrModelNotFound  = errors.New("device model not found")
	ErrGitCommit      = errors.New("git commit failed")
	ErrConfigEmpty    = errors.New("empty configuration")
	ErrDeviceNotFound = errors.New("device not found")
	ErrGroupNoDevices = errors.New("group has no devices")
	ErrSchedulerDown  = errors.New("scheduler unavailable")
)

// Error is a typed error carrying a sentinel Kind and an original human-readable
// message. Error() keeps the original message (byte-for-byte compatible with the
// previous fmt.Errorf strings), while Unwrap exposes both the Kind and the
// wrapped cause for errors.Is/As.
type Error struct {
	kind  error
	cause error
	msg   string
}

func (e *Error) Error() string { return e.msg }

// Unwrap allows errors.Is(err, ErrXxx) and errors.Is(err, underlying) to work
// through wrapping layers.
func (e *Error) Unwrap() []error {
	chain := []error{e.kind}
	if e.cause != nil {
		chain = append(chain, e.cause)
	}
	return chain
}

// Kind returns the sentinel error associated with this error.
func (e *Error) Kind() error { return e.kind }

// Wrap builds a typed error with an unformatted message.
func Wrap(kind error, msg string) error {
	if msg == "" {
		msg = kind.Error()
	}
	return &Error{kind: kind, msg: msg}
}

// Wrapf builds a typed error with a formatted message.
func Wrapf(kind error, format string, args ...interface{}) error {
	return &Error{kind: kind, msg: fmt.Sprintf(format, args...)}
}

// WrapCause builds a typed error that also preserves the underlying cause,
// equivalent to fmt.Errorf("<kind marker>: "+format, args...).
func WrapCause(kind error, cause error, format string, args ...interface{}) error {
	return &Error{kind: kind, cause: cause, msg: fmt.Sprintf(format, args...)}
}
