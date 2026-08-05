package cli

import "errors"

// ReportedError means the command already wrote a user-facing message.
// Callers should exit non-zero without printing the error again.
type ReportedError struct {
	Err error
}

func (e *ReportedError) Error() string {
	if e.Err == nil {
		return "command failed"
	}
	return e.Err.Error()
}

func (e *ReportedError) Unwrap() error {
	return e.Err
}

func reported(err error) error {
	if err == nil {
		return &ReportedError{}
	}
	return &ReportedError{Err: err}
}

func IsReported(err error) bool {
	var r *ReportedError
	return errors.As(err, &r)
}
