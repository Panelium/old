package errors

import "errors"

var UserNotFound error = errors.New("user not found")
var InvalidCredentials error = errors.New("invalid credentials")
var SessionCreationFailed error = errors.New("session creation failed")
var TokenCreationFailed error = errors.New("token creation failed")
