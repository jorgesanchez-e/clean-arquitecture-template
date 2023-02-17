package interfaceadapters

import (
	"context"

	"clean-arquitecture-template/internal/domain/register"
	"clean-arquitecture-template/internal/interfaceadapters/storage/memory"
	"clean-arquitecture-template/internal/interfaceadapters/storage/mongodb"
)

type MemRepoService struct {
	memRepo register.LineRepository
}

func NewMemRepoService(ctx context.Context) MemRepoService {
	return MemRepoService{
		memRepo: memory.New(ctx),
	}
}

type MongoRepoService struct {
	mongoRepo register.LineRepository
}

func NewMongoRepoService(ctx context.Context, cnf mongodb.Config) MongoRepoService {
	return MongoRepoService{
		mongoRepo: mongodb.New(ctx, cnf),
	}
}
