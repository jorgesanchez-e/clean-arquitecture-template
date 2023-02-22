package queries

import (
	"clean-arquitecture-template/internal/domain/example"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewGetExampleRequestHandlerHandle(t *testing.T) {
	ctx := context.Background()
	tstamp := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60)).UTC()
	data := "first-line"

	newID := "hello"

	type fields struct {
		repo     example.LineRepository
		provider example.IdentityProvider
	}

	type args struct {
		req GetExampleRequest
	}

	testCases := []struct {
		testName       string
		fields         fields
		args           args
		expectedResult *GetExampleResult
		expectedError  error
	}{
		{
			testName: "parsing-id-error-case",
			fields: fields{
				repo: func() example.MockRepository {
					return example.MockRepository{}
				}(),

				provider: func() example.MockIdentityProvider {
					idProvider := example.MockIdentityProvider{}
					idProvider.On("ParseID", newID).Return(example.MockIdentifier(""), errors.New("invalid-id"))

					return idProvider
				}(),
			},
			args: args{
				req: GetExampleRequest{
					ID: newID,
				},
			},
			expectedError: ErrInvalidID,
		},
		{
			testName: "system-error-case",
			fields: fields{
				repo: func() example.MockRepository {
					mr := example.MockRepository{}
					mr.On("Read", ctx, example.MockIdentifier(newID)).Return(&example.Line{}, errors.New("some-error"))
					return mr
				}(),

				provider: func() example.MockIdentityProvider {
					idProvider := example.MockIdentityProvider{}
					idProvider.On("ParseID", newID).Return(example.MockIdentifier(newID), nil)

					return idProvider
				}(),
			},
			args: args{
				req: GetExampleRequest{
					ID: newID,
				},
			},
			expectedError: ErrSystem,
		},
		{
			testName: "success-case",
			fields: fields{
				repo: func() example.MockRepository {
					mr := example.MockRepository{}

					mr.On("Read", ctx, example.MockIdentifier(newID)).Return(&example.Line{
						ID:      example.MockIdentifier(newID),
						Created: tstamp,
						Data:    data,
					}, nil)

					return mr
				}(),
				provider: func() example.MockIdentityProvider {
					idProvider := example.MockIdentityProvider{}
					idProvider.On("ParseID", newID).Return(example.MockIdentifier(newID), nil)

					return idProvider
				}(),
			},
			args: args{
				req: GetExampleRequest{
					ID: newID,
				},
			},
			expectedResult: &GetExampleResult{
				ID:        newID,
				Data:      data,
				CreatedAt: tstamp,
			},
			expectedError: nil,
		},
	}

	for _, c := range testCases {
		name := c.testName
		repo := c.fields.repo
		provider := c.fields.provider
		request := c.args.req
		expectedError := c.expectedError
		expectedResult := c.expectedResult

		t.Run(name, func(t *testing.T) {
			h := NewGetExampleRequestHandler(repo, provider)
			result, err := h.Handle(ctx, request)

			assert.Equal(t, expectedResult, result)
			assert.ErrorIs(t, err, expectedError)
		})
	}

}
