package example

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mr MockRepository) Write(ctx context.Context, line Line) error {
	args := mr.Called(ctx, line)
	return args.Error(0)
}

func (mr MockRepository) Read(ctx context.Context, id Identifier) (*Line, error) {
	args := mr.Called(ctx, id)
	return args.Get(0).(*Line), args.Error(1)
}

type MockIdentityProvider struct {
	mock.Mock
}

func (mip MockIdentityProvider) NewID() Identifier {
	args := mip.Called()

	return args.Get(0).(Identifier)
}

func (mip MockIdentityProvider) ParseID(ids string) (Identifier, error) {
	args := mip.Called(ids)

	return args.Get(0).(Identifier), args.Error(1)
}
