package s3manager

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Errors that may be returned from an S3 client.
var (
	ErrBucketDoesNotExist = errors.New("The specified bucket does not exist.") // nolint: golint
	ErrKeyDoesNotExist    = errors.New("The specified key does not exist.")    // nolint: golint
)

// handleHTTPError handles HTTP errors.
func handleHTTPError(w http.ResponseWriter, err error) {
	var code int

	switch err := errors.Cause(err).(type) {
	case *json.SyntaxError:
		code = http.StatusUnprocessableEntity
	default:
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			code = http.StatusUnprocessableEntity
		} else if err == ErrBucketDoesNotExist || err == ErrKeyDoesNotExist {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}
	}

	http.Error(w, http.StatusText(code), code)

	// Log if server error
	if code >= 500 {
		log.Println(err)
	}
}
