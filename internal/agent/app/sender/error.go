package sender

import "fmt"

type ErrSendRequest struct {
	Err error
}

func NewErrSendRequest(err error) *ErrSendRequest {
	return &ErrSendRequest{Err: err}
}

func (e *ErrSendRequest) Error() string {
	return fmt.Sprintf("error sending request: %s", e.Err)
}

func (e *ErrSendRequest) Unwrap() error {
	return e.Err
}
