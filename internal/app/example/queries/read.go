package queries

import (
	"context"
	"time"

	"clean-arquitecture-template/internal/domain/example"
)

// Read(context.Context, Identifier) (*Line, error)
type GetExampleRequest struct {
	ID string
}

type GetExampleResult struct {
	ID        string
	Data      string
	CreatedAt time.Time
}

type GetExampleRequestHandler interface {
	Handle(ctx context.Context, req GetExampleRequest) (*GetExampleResult, error)
}

type getExampleRequestHandler struct {
	repo       example.LineRepository
	idProvider example.IdentityProvider
}

func NewGetExampleRequestHandler(repo example.LineRepository, idProvider example.IdentityProvider) GetExampleRequestHandler {
	return getExampleRequestHandler{
		repo:       repo,
		idProvider: idProvider,
	}
}

func (h getExampleRequestHandler) Handle(ctx context.Context, req GetExampleRequest) (*GetExampleResult, error) {
	id, err := h.idProvider.ParseID(req.ID)
	if err != nil {
		return nil, err
	}

	line, err := h.repo.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetExampleResult{
		ID:        line.ID.String(),
		CreatedAt: line.Created,
		Data:      line.Data,
	}, nil
}
