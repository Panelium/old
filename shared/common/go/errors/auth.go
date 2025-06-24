package errors

import "errors"

var UserNotFound = errors.New("user not found")
var InvalidCredentials = errors.New("invalid credentials")
var SessionCreationFailed = errors.New("session creation failed")
