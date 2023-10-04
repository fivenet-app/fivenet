package common

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ErrorKeyFormat = "errors.%s.%s"

type Error struct {
	Title   string
	Content string
}

type Errors struct {
	service string
}

func (e *Errors) Error(code codes.Code, name string) error {
	return status.Error(code, fmt.Sprintf(ErrorKeyFormat, e.service, name))
}
