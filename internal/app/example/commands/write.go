package commands

import (
	"context"

	"clean-arquitecture-template/internal/domain/example"
)

type AddExampleRequest struct {
	Data string
}

type CreateLineRequestHandler interface {
	Handle(ctx context.Context, command AddExampleRequest) (*string, error)
}

type addExampleRequestHandler struct {
	repo       example.LineRepository
	idProvider example.IdentityProvider
}

func NewAddExampleRequestHandler(repo example.LineRepository, idProvider example.IdentityProvider) CreateLineRequestHandler {
	return addExampleRequestHandler{
		repo:       repo,
		idProvider: idProvider,
	}
}

func (h addExampleRequestHandler) Handle(ctx context.Context, command AddExampleRequest) (*string, error) {
	line := example.Line{
		ID: h.idProvider.NewID(),
	}

	err := h.repo.Write(ctx, line)
	if err != nil {
		return nil, err
	}

	id := line.ID.String()

	return &id, nil
}
