package main

import (
	"context"

	exampleapp "clean-arquitecture-template/internal/app/example"
	"clean-arquitecture-template/internal/inputports/example"
	"clean-arquitecture-template/internal/interfaceadapters/storage/mongodb"
)

var cnf config

func main() {
	ctx := context.Background()

	repo := mongodb.NewExampleRepo(ctx, cnf)
	idProv := mongodb.NewIdentityProvider()
	services := exampleapp.NewServices(repo, idProv)
	rest := example.NewServices(ctx, services)

	rest.Server.ListenAndServe("8080")
}

type config struct{}

func (c config) GetDSN() string {
	return ""
}

func (c config) DatabaseName() string {
	return ""
}

func (c config) TableName() string {
	return ""
}
