package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	app "clean-arquitecture-template/internal/app/example"
	"clean-arquitecture-template/internal/app/example/commands"
	"clean-arquitecture-template/internal/app/example/queries"
)

type mockCommandCreateLineHandler struct {
	Handler func(context.Context, commands.AddExampleRequest) (*string, error)
}

func (m mockCommandCreateLineHandler) Handle(ctx context.Context, command commands.AddExampleRequest) (*string, error) {
	return m.Handler(ctx, command)
}

func Test_WriteAppExampleHandler(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		testName         string
		handler          commands.CreateLineRequestHandler
		requestData      string
		expectedHTTPCode int
		expectedResponse string
	}{
		{
			testName: "system-error-test",
			handler: mockCommandCreateLineHandler{Handler: func(ctx context.Context, req commands.AddExampleRequest) (*string, error) {
				return nil, commands.ErrSystem
			}},
			expectedHTTPCode: 500,
			expectedResponse: "\"system error\"\n",
		},
		{
			testName: "unknown-error-test",
			handler: mockCommandCreateLineHandler{Handler: func(ctx context.Context, req commands.AddExampleRequest) (*string, error) {
				return nil, nil
			}},
			expectedHTTPCode: 500,
			expectedResponse: "\"system error\"\n",
		},

		{
			testName: "success-test",
			handler: mockCommandCreateLineHandler{Handler: func(ctx context.Context, req commands.AddExampleRequest) (*string, error) {
				newID := "1234567890"
				return &newID, nil
			}},
			requestData:      `{"data":"x"}`,
			expectedHTTPCode: 200,
			expectedResponse: "{\n \"new_id\": \"1234567890\"\n}\n",
		},
		{
			testName: "success-test",
			handler: mockCommandCreateLineHandler{Handler: func(ctx context.Context, req commands.AddExampleRequest) (*string, error) {
				newID := "1234567890"
				return &newID, nil
			}},
			requestData:      `{"data":"x"`,
			expectedHTTPCode: 400,
			expectedResponse: "code=400, message=code=400, message=unexpected EOF, internal=unexpected EOF: input param error",
		},
	}

	for _, c := range testCases {
		testName := c.testName
		handler := c.handler
		expectedCode := c.expectedHTTPCode
		expectedResponse := c.expectedResponse

		server := Server{
			ctx:    ctx,
			server: echo.New(),
			exampleServices: app.Services{
				ExampleService: app.ExampleServices{
					Commands: app.Commands{
						CreateExampleHandler: handler,
					},
				},
			},
		}

		req := httptest.NewRequest(http.MethodPost, "/example/write", strings.NewReader(c.requestData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		t.Run(testName, func(t *testing.T) {
			c := server.server.NewContext(req, rec)

			err := server.writeAppExample(c)
			code := rec.Code
			response := rec.Body.String()

			if err != nil {
				he, ok := err.(*echo.HTTPError)
				if ok {
					assert.Equal(t, expectedCode, he.Code)
					assert.Equal(t, expectedResponse, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedCode, code)
				assert.Equal(t, expectedResponse, response)
			}
		})
	}
}

type mockCommandReadLineHandler struct {
	Handler func(context.Context, queries.GetExampleRequest) (*queries.GetExampleResult, error)
}

func (m mockCommandReadLineHandler) Handle(ctx context.Context, req queries.GetExampleRequest) (*queries.GetExampleResult, error) {
	return m.Handler(ctx, req)
}

func Test_ReadAppExample(t *testing.T) {
	const requestParamName string = "id"

	ctx := context.Background()
	tstamp := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60)).UTC()

	testCases := []struct {
		testName         string
		handler          queries.GetExampleRequestHandler
		requestPath      string
		requestValue     string
		expectedHTTPCode int
		expectedResponse string
	}{
		{
			testName: "system-error-test",
			handler: mockCommandReadLineHandler{Handler: func(ctx context.Context, req queries.GetExampleRequest) (*queries.GetExampleResult, error) {
				return nil, commands.ErrSystem
			}},
			requestPath:      "/example/read/1000",
			requestValue:     "1000",
			expectedHTTPCode: 500,
			expectedResponse: "\"system error\"\n",
		},
		{
			testName: "success-test",
			handler: mockCommandReadLineHandler{Handler: func(ctx context.Context, req queries.GetExampleRequest) (*queries.GetExampleResult, error) {
				return &queries.GetExampleResult{
					ID:        req.ID,
					Data:      "first-line",
					CreatedAt: tstamp,
				}, nil
			}},
			requestPath:      "/example/read/1000",
			requestValue:     "1000",
			expectedHTTPCode: 200,
			expectedResponse: "{\n \"id\": \"1000\",\n \"created_at\": \"2018-09-16 10:00:00 +0000 UTC\",\n \"data\": \"first-line\"\n}\n",
		},
	}

	for _, c := range testCases {
		testName := c.testName
		handler := c.handler
		expectedCode := c.expectedHTTPCode
		expectedResponse := c.expectedResponse
		requestPath := c.requestPath
		requestValue := c.requestValue

		server := Server{
			ctx:    ctx,
			server: echo.New(),
			exampleServices: app.Services{
				ExampleService: app.ExampleServices{
					Queries: app.Queries{
						ReadExampleHandler: handler,
					},
				},
			},
		}

		req := httptest.NewRequest(http.MethodGet, requestPath, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		t.Run(testName, func(t *testing.T) {
			c := server.server.NewContext(req, rec)
			c.SetPath(requestPath)
			c.SetParamNames(requestParamName)
			c.SetParamValues(requestValue)

			err := server.readAppExample(c)
			code := rec.Code
			response := rec.Body.String()

			assert.NoError(t, err)
			assert.Equal(t, expectedCode, code)
			assert.Equal(t, expectedResponse, response)
		})
	}
}
