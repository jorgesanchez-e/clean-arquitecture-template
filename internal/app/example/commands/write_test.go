package commands

import (
	"clean-arquitecture-template/internal/domain/example"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddExampleRequestHandlerHandle(t *testing.T) {
	ctx := context.Background()
	data := "first-line"

	newID := "hello"

	type fields struct {
		repo       example.LineRepository
		idProvider example.IdentityProvider
	}

	type args struct {
		request AddExampleRequest
	}

	testCases := []struct {
		name          string
		fields        fields
		args          args
		expectedNewID *string
		expectedError error
	}{
		{
			name: "successfull-case",
			fields: fields{
				repo: func() example.MockRepository {
					mr := example.MockRepository{}

					mr.On("Write", ctx, example.Line{
						ID:   example.MockIdentifier(newID),
						Data: data,
					}).Return(nil)

					return mr
				}(),
				idProvider: func() example.MockIdentityProvider {
					provider := example.MockIdentityProvider{}
					provider.On("NewID").Return(example.MockIdentifier(newID))

					return provider
				}(),
			},
			args: args{
				request: AddExampleRequest{
					Data: "first-line",
				},
			},
			expectedNewID: &newID,
			expectedError: nil,
		},
		{
			name: "error-case",
			fields: fields{
				repo: func() example.MockRepository {
					mr := example.MockRepository{}

					mr.On("Write", ctx, example.Line{
						ID:   example.MockIdentifier(newID),
						Data: data,
					}).Return(errors.New("some-error"))

					return mr
				}(),
				idProvider: func() example.MockIdentityProvider {
					provider := example.MockIdentityProvider{}
					provider.On("NewID").Return(example.MockIdentifier(newID))

					return provider
				}(),
			},
			args: args{
				request: AddExampleRequest{
					Data: "first-line",
				},
			},
			expectedNewID: nil,
			expectedError: ErrSystem,
		},
	}

	for _, c := range testCases {
		name := c.name
		repo := c.fields.repo
		idProvider := c.fields.idProvider
		request := c.args.request

		expectedNewID := c.expectedNewID
		expectedError := c.expectedError

		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			h := NewAddExampleRequestHandler(repo, idProvider)
			newID, err := h.Handle(ctx, request)

			assert.Equal(t, expectedNewID, newID)
			assert.ErrorIs(t, err, expectedError)
		})
	}
}
