package http

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"

	app "clean-arquitecture-template/internal/app/example"
)

const (
	exampleRoute string = "/example"

	writePath string = "/write"
	readPath  string = "/read/:id"
)

type Server struct {
	ctx             context.Context
	exampleServices app.Services
	server          *echo.Echo
}

func NewServer(ctx context.Context, appServices app.Services) Server {
	s := Server{
		ctx:             ctx,
		exampleServices: appServices,
		server:          echo.New(),
	}

	s.initApi()

	return s
}

func (s Server) initApi() {
	g := s.server.Group(exampleRoute)

	g.POST(writePath, s.writeAppExample)
	g.GET(readPath, s.readAppExample)
}

func (s Server) ListenAndServe(port string) {
	err := s.server.Start(port)

	log.Fatal(err)
}
