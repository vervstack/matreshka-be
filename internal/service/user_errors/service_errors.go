package user_errors

import (
	errors "go.redsock.ru/rerrors"
	"google.golang.org/grpc/codes"
)

var (
	ErrValidation    = errors.NewUserError("Validation error", codes.InvalidArgument)
	ErrAlreadyExists = errors.NewUserError("Already exists", codes.AlreadyExists)
	ErrNotFound      = errors.NewUserError("Not found", codes.NotFound)

	ErrNotImplemented = errors.NewUserError("Not implemented", codes.Unimplemented)
)
