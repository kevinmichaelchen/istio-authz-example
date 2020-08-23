package grpc

import (
	"errors"
)

var ErrMalformedAuthHeader = errors.New("malformed Authorization header; must be two elements, 'Bearer' followed by token")
