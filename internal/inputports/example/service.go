package example

import (
	app "clean-arquitecture-template/internal/app/example"
	"clean-arquitecture-template/internal/inputports/example/http"
	"context"
)

type Services struct {
	Server http.Server
}

func NewServices(ctx context.Context, app app.Services) Services {
	return Services{
		Server: http.NewServer(ctx, app),
	}
}
