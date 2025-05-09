package httperror

import (
	"net/http"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

func Error(err error) (string, int) {
	return getMessage(err), getStatus(err)
}

func getStatus(err error) int {
	momoErr, ok := err.(*momoError.MomoError)

	if !ok {
		return http.StatusInternalServerError
	}
	errType := momoErr.GetErrorType()
	return mapMomoErrorTypeToHttpStatus(errType)
}

func mapMomoErrorTypeToHttpStatus(errType momoError.ErrorType) int {
	switch errType {
	case momoError.BadRequest:
		return http.StatusBadRequest
	case momoError.Forbidden:
		return http.StatusForbidden
	case momoError.UnExpected:
		return http.StatusInternalServerError
	case momoError.NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func getMessage(err error) string {
	momoErr, ok := err.(*momoError.MomoError)
	if !ok {
		return "something went wrong"
	}

	message := momoErr.Message()

	if len(message) > 0 {
		return message
	}

	code := getStatus(err)
	switch code {
	case http.StatusBadRequest:
		return "input is wrong"
	case http.StatusNotFound:
		return "no record found"
	default:
		return "something went wrong"
	}
}
