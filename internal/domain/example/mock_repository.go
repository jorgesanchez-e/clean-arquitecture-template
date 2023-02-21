package example

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
	CreatedTime time.Time
}

func (mr MockRepository) Write(ctx context.Context, line Line) error {
	line.Created = mr.CreatedTime
	args := mr.Called(ctx, line)
	return args.Error(0)
}

func (mr MockRepository) Read(ctx context.Context, id Identifier) (*Line, error) {
	args := mr.Called(ctx, id)
	return args.Get(0).(*Line), args.Error(1)
}

type MockIdentityProvider struct {
	mock.Mock
	ID MockIdentifier
}

type MockIdentifier string

func (mid MockIdentifier) String() string {
	return string(mid)
}

func (mip MockIdentityProvider) NewID() Identifier {
	mip.Called()

	return Identifier(mip.ID)
}

func (mip MockIdentityProvider) ParseID(ids string) (Identifier, error) {
	args := mip.Called(ids)

	return args.Get(0).(Identifier), args.Error(1)
}
