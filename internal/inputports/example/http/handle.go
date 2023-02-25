package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"clean-arquitecture-template/internal/app/example/commands"
	"clean-arquitecture-template/internal/app/example/queries"
)

type WriteExampleRequest struct {
	Data string `json:"data"`
}

type WriteExampleResponse struct {
	NewID string `json:"new_id"`
}

func (s Server) writeAppExample(c echo.Context) error {
	data := new(WriteExampleRequest)

	response := NewResponser(c)

	if err := c.Bind(data); err != nil {
		response.WithHTTPError(fmt.Errorf("%s: %w", err.Error(), ErrInputParam))
	} else {
		ctx, cancel := context.WithTimeout(s.ctx, time.Second)
		defer cancel()

		id, err := s.exampleServices.ExampleService.Commands.CreateExampleHandler.Handle(ctx, commands.AddExampleRequest{Data: data.Data})
		if err != nil {
			response.WithJSONError(err)
		} else {
			if id == nil {
				response.WithJSONError(commands.ErrSystem)
			} else {
				response.WithJSON(http.StatusOK, WriteExampleResponse{
					NewID: *id,
				})
			}
		}
	}

	return response.Response()
}

type readAppExampleResponse struct {
	ID        string `json:"id"`
	CreatedAT string `json:"created_at"`
	Data      string `json:"data"`
}

func (s Server) readAppExample(c echo.Context) error {
	idParam := c.Param("id")

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	response := NewResponser(c)
	if result, err := s.exampleServices.ExampleService.Queries.ReadExampleHandler.Handle(ctx, queries.GetExampleRequest{ID: idParam}); err != nil {
		response.WithJSONError(err)
	} else {
		response.WithJSON(http.StatusOK, readAppExampleResponse{
			ID:        result.ID,
			CreatedAT: result.CreatedAt.String(),
			Data:      result.Data,
		})
	}

	return response.Response()
}
