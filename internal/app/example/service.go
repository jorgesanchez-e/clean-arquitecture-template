package example

import (
	"clean-arquitecture-template/internal/app/example/commands"
	"clean-arquitecture-template/internal/app/example/queries"
	"clean-arquitecture-template/internal/domain/example"
)

type Commands struct {
	CreateExampleHandler commands.CreateLineRequestHandler
}

type Queries struct {
	ReadExampleHandler queries.GetExampleRequestHandler
}

type ExampleServices struct {
	Commands Commands
	Queries  Queries
}

type Services struct {
	ExampleService ExampleServices
}

func NewServices(examRepo example.LineRepository, idProdiver example.IdentityProvider) Services {
	return Services{
		ExampleService: ExampleServices{
			Commands: Commands{
				CreateExampleHandler: commands.NewAddExampleRequestHandler(examRepo, idProdiver),
			},
			Queries: Queries{
				ReadExampleHandler: queries.NewGetExampleRequestHandler(examRepo, idProdiver),
			},
		},
	}
}
