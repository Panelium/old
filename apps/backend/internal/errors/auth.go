package errors

import "errors"

var UserNotFound error = errors.New("user not found")
var InvalidCredentials error = errors.New("invalid credentials")
