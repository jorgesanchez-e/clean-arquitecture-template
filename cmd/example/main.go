package main

import (
	"context"
	"log"

	"clean-arquitecture-template/config"
	exampleapp "clean-arquitecture-template/internal/app/example"
	"clean-arquitecture-template/internal/inputports/example"
	"clean-arquitecture-template/internal/inputports/example/http"
	"clean-arquitecture-template/internal/interfaceadapters/storage/mongodb"
)

func main() {
	ctx := context.Background()

	cnf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	mongoConfig, err := mongodb.ReadConfig(cnf)
	if err != nil {
		log.Fatal(err)
	}

	restConf, err := http.ReadConfig(cnf)
	if err != nil {
		log.Fatal(err)
	}

	repo := mongodb.NewExampleRepo(ctx, mongoConfig)
	idProv := mongodb.NewIdentityProvider()
	services := exampleapp.NewServices(repo, idProv)
	rest := example.NewServices(ctx, services, restConf)

	rest.Server.ListenAndServe()
}
