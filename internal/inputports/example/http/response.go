package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"clean-arquitecture-template/internal/app/example/commands"
	"clean-arquitecture-template/internal/app/example/queries"
)

const (
	httpErrorResponse responseType = iota
	jsonErrorResponse
	jsonResponse
	notFoundResponse

	defaultResponserError string = "response not set"
	defaultResponserCode  int    = http.StatusNotImplemented

	ErrInputParam Error = "input param error"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type responseType int

type responser struct {
	echoContext  echo.Context
	responseType responseType
	code         int
	payload      interface{}
}

func NewResponser(c echo.Context) *responser {
	return &responser{
		responseType: -1,
		echoContext:  c,
	}
}

func (r *responser) withError(err error) *responser {
	if r == nil {
		return r
	}

	if errors.Is(err, commands.ErrSystem) || errors.Is(err, queries.ErrSystem) {
		r.code = http.StatusInternalServerError
	}

	if errors.Is(err, queries.ErrInvalidID) || errors.Is(err, ErrInputParam) {
		r.code = http.StatusBadRequest
	}

	r.payload = err.Error()

	return r
}

func (r *responser) WithHTTPError(err error) *responser {
	if r == nil {
		return r
	}

	r.withError(err)
	r.responseType = httpErrorResponse

	return r
}

func (r *responser) WithJSONError(err error) *responser {
	if r == nil {
		return r
	}

	r.withError(err)
	r.responseType = jsonErrorResponse

	return r
}

func (r *responser) WithJSON(code int, payload interface{}) *responser {
	if r == nil {
		return r
	}

	r.responseType = jsonResponse
	r.code = code
	r.payload = payload

	return r
}

func (r *responser) WithNotFound() *responser {
	if r == nil {
		return r
	}

	r.code = http.StatusNotFound
	r.responseType = notFoundResponse

	return r
}

func (r *responser) Response() error {
	if r == nil {
		log.Error(defaultResponserError)
		return nil
	}

	switch r.responseType {
	case httpErrorResponse:
		return echo.NewHTTPError(r.code, r.payload)
	case jsonErrorResponse, jsonResponse:
		return r.echoContext.JSONPretty(r.code, r.payload, " ")
	case notFoundResponse:
		return r.echoContext.NoContent(r.code)
	}

	return echo.NewHTTPError(defaultResponserCode, defaultResponserError)
}
