package exceptions

import (
	"net/http"
)

const (
	ERRDOMAIN       Status = 1
	ERRBUSSINESS    Status = 2
	ERRSYSTEM       Status = 3
	ERRNOTFOUND     Status = 4
	ERRREPOSITORY   Status = 5
	ERRUNKNOWN      Status = 6
	ERRAUTHORIZED   Status = 7
	ERRFORBIDDEN    Status = 8
	ERRDATACONFLICT Status = 9
)

type Status int

type CustomError struct {
	Status Status
	Errors error
}

func MapToHTTPStatusCode(status Status) int {
	var httpStatusCode int
	switch status {
	case ERRDOMAIN:
		httpStatusCode = http.StatusBadRequest
	case ERRBUSSINESS:
		httpStatusCode = http.StatusBadRequest
	case ERRSYSTEM:
		httpStatusCode = http.StatusInternalServerError
	case ERRNOTFOUND:
		httpStatusCode = http.StatusNotFound
	case ERRREPOSITORY:
		httpStatusCode = http.StatusInternalServerError
	case ERRUNKNOWN:
		httpStatusCode = http.StatusInternalServerError
	case ERRAUTHORIZED:
		httpStatusCode = http.StatusUnauthorized
	case ERRFORBIDDEN:
		httpStatusCode = http.StatusForbidden
	case ERRDATACONFLICT:
		httpStatusCode = http.StatusConflict
	default:
		httpStatusCode = http.StatusInternalServerError
	}

	return httpStatusCode
}
