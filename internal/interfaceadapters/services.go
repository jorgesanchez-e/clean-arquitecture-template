package interfaceadapters

import (
	"context"

	"clean-arquitecture-template/internal/domain/example"
	"clean-arquitecture-template/internal/interfaceadapters/example/storage/memory"
	"clean-arquitecture-template/internal/interfaceadapters/example/storage/mongodb"
)

type MemExampleRepoService struct {
	memRepo example.LineRepository
}

func NewMemRepoService(ctx context.Context) MemExampleRepoService {
	return MemExampleRepoService{
		memRepo: memory.NewExampleRepo(ctx),
	}
}

type MongoExampleRepoService struct {
	mongoRepo example.LineRepository
}

func NewMongoRepoService(ctx context.Context, cnf mongodb.Config) MongoExampleRepoService {
	return MongoExampleRepoService{
		mongoRepo: mongodb.NewExampleRepo(ctx, cnf),
	}
}
