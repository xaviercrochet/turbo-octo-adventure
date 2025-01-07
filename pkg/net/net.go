package net

import (
	"errors"
	"net/http"
)

var (
	ErrNoAccess         = errors.New("not authorized")
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrNotFound         = errors.New("resource not found")
	ErrGeneric          = errors.New("request failed")
)

/*

Return an error based on the http status code of the response

*/

func HttpStatusCodeToErr(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return ErrNotAuthenticated
	case http.StatusForbidden:
		return ErrNoAccess
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return ErrGeneric
	}
}
