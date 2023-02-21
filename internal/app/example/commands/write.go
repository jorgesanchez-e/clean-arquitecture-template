package commands

import (
	"context"
	"fmt"

	"clean-arquitecture-template/internal/domain/example"
)

const (
	ErrSystem ServiceError = "system error"
)

type ServiceError string

func (se ServiceError) Error() string {
	return string(se)
}

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
		ID:   h.idProvider.NewID(),
		Data: command.Data,
	}

	err := h.repo.Write(ctx, line)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrSystem)
	}

	id := line.ID.String()

	return &id, nil
}
