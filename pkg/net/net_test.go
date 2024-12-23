package net

import (
	"net/http"
	"testing"
)

func TestHttpStatusCodeToErr(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
	}{
		{
			name:       "should return nil for StatusOK",
			statusCode: http.StatusOK,
			err:        nil,
		},
		{
			name:       "should return ErrUnauthorized for StatusUnauthorized",
			statusCode: http.StatusUnauthorized,
			err:        ErrUnauthorized,
		},
		{
			name:       "should return ErrNoAccess for StatusForbidden",
			statusCode: http.StatusForbidden,
			err:        ErrNoAccess,
		},
		{
			name:       "should return ErrNotFound for StatusNotFound",
			statusCode: http.StatusNotFound,
			err:        ErrNotFound,
		},
		{
			name:       "should return ErrGeneric for unknown status code",
			statusCode: http.StatusInternalServerError,
			err:        ErrGeneric,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
			}

			err := HttpStatusCodeToErr(resp)

			if err != tt.err {
				t.Errorf("HttpStatusCodeToErr() error = %v, err %v", err, tt.err)
			}
		})
	}
}
