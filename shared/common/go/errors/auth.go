package errors

import (
	"connectrpc.com/connect"
	"errors"
)

var ConnectInvalidCredentials = connect.NewError(connect.CodeUnauthenticated, InvalidCredentials)

var UserNotFound = errors.New("user not found")
var InvalidCredentials = errors.New("invalid credentials")
var SessionCreationFailed = errors.New("session creation failed")
var SessionNotFound = errors.New("session not found")
