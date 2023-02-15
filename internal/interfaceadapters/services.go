package interfaceadapters

import (
	"context"

	"clean-arquitecture-template/internal/domain/register"
	"clean-arquitecture-template/internal/interfaceadapters/storage/memory"
)

type Services struct {
	LineRepository register.LineRepository
}

func NewServices(ctx context.Context) Services {
	return Services{
		LineRepository: memory.New(ctx),
	}
}
