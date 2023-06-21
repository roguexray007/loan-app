package utils

import (
	goerr "errors"
)

// GetError Gets error interface given the interface of error
// Extract the error data from the err type
// If type couldn't be determined then create a default error
func GetError(err interface{}) error {
	var er error

	switch err.(type) {
	case error:
		er = err.(error)
	case string:
		er = goerr.New(err.(string))
	default:
		er = goerr.New("internal server error")
	}

	return er
}
