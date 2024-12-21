package net

import (
	"errors"
	"net/http"
)

var (
	ErrNoAccess     = errors.New("authentication failed")
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("resource not found")
	ErrGeneric      = errors.New("request failed")
)

/*

Return an error based on the http status code of the response

*/

func HttpStatusCodeToErr(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrNoAccess
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return ErrGeneric
	}
}
