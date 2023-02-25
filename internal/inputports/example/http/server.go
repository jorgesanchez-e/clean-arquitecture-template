package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/labstack/echo/v4"

	app "clean-arquitecture-template/internal/app/example"
)

const (
	exampleRoute string = "/example"

	writePath string = "/write"
	readPath  string = "/read/:id"

	configPath string = "apps.example.input-ports.rest"

	ErrReadConfig err = "unable to read config"
)

type err string

func (e err) Error() string {
	return string(e)
}

type Config interface {
	Address() string
}

type config struct {
	Addr string `json:"address"`
	Port string `json:"port"`
}

func (cnf config) Address() string {
	return fmt.Sprintf("%s:%s", cnf.Addr, cnf.Port)
}

type ConfigReader interface {
	Find(node string) (io.Reader, error)
}

func ReadConfig(cnfr ConfigReader) (Config, error) {
	reader, err := cnfr.Find(configPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrReadConfig)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrReadConfig)
	}

	cnf := config{}
	if err = json.Unmarshal(data, &cnf); err != nil {
		return nil, ErrReadConfig
	}

	return cnf, nil
}

type Server struct {
	ctx             context.Context
	exampleServices app.Services
	server          *echo.Echo
	address         string
}

func NewServer(ctx context.Context, appServices app.Services, cnf Config) Server {
	s := Server{
		ctx:             ctx,
		exampleServices: appServices,
		server:          echo.New(),
		address:         cnf.Address(),
	}

	s.initApi()

	return s
}

func (s Server) initApi() {
	g := s.server.Group(exampleRoute)

	g.POST(writePath, s.writeAppExample)
	g.GET(readPath, s.readAppExample)
}

func (s Server) ListenAndServe() {
	err := s.server.Start(s.address)

	log.Fatal(err)
}
