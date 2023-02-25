package http

import (
	"clean-arquitecture-template/internal/app/example/commands"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_Responser(t *testing.T) {
	econtext := echo.New().NewContext(nil, nil)

	testCases := []struct {
		testName          string
		expectedResponser *responser
		buildResponser    func() *responser
	}{
		{
			testName:          "nil-case",
			expectedResponser: nil,
			buildResponser: func() *responser {
				var r *responser

				r.withError(nil)
				r.WithHTTPError(nil)
				r.WithJSONError(nil)
				r.WithJSON(200, nil)
				r.WithNotFound()
				r.Response()

				return r
			},
		},
		{
			testName: "http-error-case",
			expectedResponser: &responser{
				echoContext:  econtext,
				responseType: httpErrorResponse,
				payload:      "some-error",
			},
			buildResponser: func() *responser {
				resp := &responser{
					echoContext: econtext,
				}

				resp = resp.WithHTTPError(errors.New("some-error"))

				return resp
			},
		},
		{
			testName: "http-system-error-case",
			expectedResponser: &responser{
				echoContext:  econtext,
				responseType: httpErrorResponse,
				code:         http.StatusInternalServerError,
				payload:      commands.ErrSystem.Error(),
			},
			buildResponser: func() *responser {
				resp := &responser{
					echoContext: econtext,
				}

				resp = resp.WithHTTPError(commands.ErrSystem)

				return resp
			},
		},
		{
			testName: "http-bad-request-error-case",
			expectedResponser: &responser{
				echoContext:  econtext,
				responseType: httpErrorResponse,
				code:         http.StatusBadRequest,
				payload:      ErrInputParam.Error(),
			},
			buildResponser: func() *responser {
				resp := &responser{
					echoContext: econtext,
				}

				resp = resp.WithHTTPError(ErrInputParam)

				return resp
			},
		},
		{
			testName: "json-error-case",
			expectedResponser: &responser{
				echoContext:  econtext,
				responseType: jsonErrorResponse,
				code:         http.StatusBadRequest,
				payload:      ErrInputParam.Error(),
			},
			buildResponser: func() *responser {
				resp := &responser{
					echoContext: econtext,
				}

				resp = resp.WithJSONError(ErrInputParam)

				return resp
			},
		},
		{
			testName: "json-case",
			expectedResponser: &responser{
				echoContext:  econtext,
				responseType: jsonResponse,
				code:         http.StatusOK,
				payload:      `{}`,
			},
			buildResponser: func() *responser {
				resp := &responser{
					echoContext: econtext,
				}

				resp = resp.WithJSON(http.StatusOK, `{}`)

				return resp
			},
		},
		{
			testName: "notfound-case",
			expectedResponser: &responser{
				echoContext:  econtext,
				responseType: notFoundResponse,
				code:         http.StatusNotFound,
				payload:      nil,
			},
			buildResponser: func() *responser {
				resp := &responser{
					echoContext: econtext,
				}

				resp = resp.WithNotFound()

				return resp
			},
		},
	}

	for _, c := range testCases {
		name := c.testName
		expectedResult := c.expectedResponser
		newResponser := c.buildResponser

		t.Run(name, func(t *testing.T) {
			result := newResponser()

			assert.Equal(t, expectedResult, result)
		})
	}

}

func Test_ResponserResponse(t *testing.T) {
	testCases := []struct {
		testName         string
		buildResponser   func() (*responser, *httptest.ResponseRecorder)
		expectedHTTPCode int
		expectedResponse string
		expectedError    error
	}{
		{
			testName: "http-error-case",
			buildResponser: func() (*responser, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()

				return NewResponser(echo.New().NewContext(nil, rec)).WithHTTPError(ErrInputParam), rec
			},
			expectedHTTPCode: 200,
			expectedError:    &echo.HTTPError{Code: 400, Message: "input param error", Internal: error(nil)},
		},
		{
			testName: "json-case",
			buildResponser: func() (*responser, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()

				r := struct {
					ID string `json:"id"`
				}{
					ID: "123",
				}

				return NewResponser(echo.New().NewContext(nil, rec)).WithJSON(http.StatusOK, r), rec
			},
			expectedHTTPCode: 200,
			expectedError:    nil,
			expectedResponse: "{\n \"id\": \"123\"\n}\n",
		},
		{
			testName: "not-found-case",
			buildResponser: func() (*responser, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()

				return NewResponser(echo.New().NewContext(nil, rec)).WithNotFound(), rec
			},
			expectedHTTPCode: 404,
			expectedError:    nil,
		},
		{
			testName: "default-error-case",
			buildResponser: func() (*responser, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()
				resp := NewResponser(echo.New().NewContext(nil, rec))

				return resp, rec
			},
			expectedHTTPCode: 200,
			expectedError:    &echo.HTTPError{Code: 501, Message: "response not set", Internal: error(nil)},
		},
	}

	for _, c := range testCases {
		name := c.testName
		responser, recorder := c.buildResponser()
		expectedHTTPCode := c.expectedHTTPCode
		expectedResponse := c.expectedResponse
		expectedError := c.expectedError

		t.Run(name, func(t *testing.T) {
			err := responser.Response()

			assert.Equal(t, expectedError, err)
			assert.Equal(t, expectedHTTPCode, recorder.Code)
			assert.Equal(t, expectedResponse, recorder.Body.String())
		})
	}
}

/*
	rec := httptest.NewRecorder()
	t.Run(testName, func(t *testing.T) {
		c := server.server.NewContext(req, rec)

*/
