package mongodb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"clean-arquitecture-template/internal/domain/example"
)

func Test_NewID(t *testing.T) {
	id := NewIdentityProvider().NewID()

	mid, is := id.(Identifier)
	assert.True(t, is)

	oid := mid.GetObjectID()
	assert.NotEmpty(t, id.String())
	assert.NotEmpty(t, mid.String())
	assert.NotEmpty(t, oid.String())
}

// func IdentifierFromText(key string) (example.Identifier, error) {
func Test_IdentifierFromText(t *testing.T) {
	testCases := []struct {
		name          string
		key           string
		identifier    Identifier
		expectedError error
	}{
		{
			name:          "successful-case",
			key:           "63f39e5fb192144de70dfa58",
			identifier:    Identifier{},
			expectedError: nil,
		},
		{
			name:          "error-case",
			key:           "hola",
			identifier:    Identifier{},
			expectedError: ErrIdentifyer,
		},
	}

	for _, c := range testCases {
		testname := c.name
		key := c.key
		id := c.identifier
		expectedError := c.expectedError

		t.Run(testname, func(t *testing.T) {
			result, err := id.ParseID(key)

			if err != nil {
				assert.ErrorIs(t, err, expectedError)
			} else {
				assert.Equal(t, true, ((example.Identifier)(result) == result))
				assert.Equal(t, 24, len(result.String()))
			}
		})
	}
}

func Test_RegisterLine(t *testing.T) {
	tstamp := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60)).UTC()
	id := primitive.NewObjectID()

	testCases := []struct {
		name           string
		input          *line
		expectedOutput *example.Line
	}{
		{
			name: "normal-case",
			input: &line{
				ID:        id,
				CreatedAT: tstamp,
				Data:      "first-line",
			},
			expectedOutput: &example.Line{
				ID:      Identifier(id),
				Created: tstamp,
				Data:    "first-line",
			},
		},
		{
			name:           "nil-case",
			input:          nil,
			expectedOutput: nil,
		},
	}

	for _, c := range testCases {
		name := c.name
		input := c.input
		expectedOutput := c.expectedOutput

		t.Run(name, func(t *testing.T) {
			result := input.registerLine()

			if result != nil {
				assert.NotEmpty(t, result.ID)
				assert.Equal(t, expectedOutput.ID.String(), result.ID.String())
				assert.Equal(t, expectedOutput.Created.String(), result.Created.String())
				assert.Equal(t, expectedOutput.Data, result.Data)
			}
		})
	}
}

func Test_Write(t *testing.T) {
	id := Identifier(primitive.NewObjectID())
	tstamp := time.Now()

	testCases := []struct {
		name          string
		input         example.Line
		ctx           context.Context
		mongoRes      bson.D
		expectedError error
	}{
		{
			name: "success-case",
			input: example.Line{
				ID:      id,
				Created: tstamp,
				Data:    "first-line",
			},
			mongoRes:      mtest.CreateSuccessResponse(),
			expectedError: nil,
		},
		{
			name: "error-id-case",
			input: example.Line{
				ID:      nil,
				Created: tstamp,
				Data:    "first-line",
			},
			mongoRes:      mtest.CreateSuccessResponse(),
			expectedError: ErrIdentifyer,
		},
		{
			name: "mogo-write-error-case",
			input: example.Line{
				ID:      id,
				Created: tstamp,
				Data:    "first-line",
			},
			mongoRes: mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   1,
				Code:    0,
				Message: "insert-one-error",
			}),
			expectedError: ErrDataInserted,
		},
	}

	for _, c := range testCases {
		testName := c.name
		ctx := c.ctx
		input := c.input
		expectedError := c.expectedError
		mongoRes := c.mongoRes

		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		defer mt.Close()

		mt.Run(testName, func(mt *mtest.T) {
			mt.AddMockResponses(mongoRes)

			st := store{
				ctx:        ctx,
				collection: mt.Coll,
			}

			err := st.Write(ctx, input)

			if err != nil {
				assert.ErrorIs(t, err, expectedError)
			} else {
				assert.Equal(t, err, expectedError)
			}
		})
	}
}

func Test_Read(t *testing.T) {
	id := Identifier(primitive.NewObjectID())
	tstamp := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60)).UTC()

	testCases := []struct {
		testName       string
		ctx            context.Context
		id             example.Identifier
		expectedResult *example.Line
		input          *line
		expectedError  error
		prepMongoMock  func(mt *mtest.T, l *line)
	}{
		{
			testName:       "identifier-and-nilctx-error-case",
			ctx:            nil,
			id:             nil,
			input:          nil,
			expectedResult: nil,
			expectedError:  ErrIdentifyer,
		},
		{
			testName:       "identifier-error-case",
			ctx:            context.Background(),
			id:             nil,
			input:          nil,
			expectedResult: nil,
			expectedError:  ErrIdentifyer,
		},
		{
			testName:       "item-not-found-case",
			ctx:            context.Background(),
			id:             id,
			input:          nil,
			expectedResult: nil,
			expectedError:  nil,
			prepMongoMock: func(mt *mtest.T, l *line) {
				ns := fmt.Sprintf("%s.%s", "dbname", "lines")
				cursorResponse := mtest.CreateCursorResponse(
					1,
					ns,
					mtest.FirstBatch)

				cursorEnd := mtest.CreateCursorResponse(
					0,
					ns,
					mtest.NextBatch)

				mt.AddMockResponses(cursorResponse, cursorEnd)
			},
		},
		{
			testName:       "mongodb-error-case",
			ctx:            context.Background(),
			id:             id,
			input:          nil,
			expectedResult: nil,
			expectedError:  ErrMongoSystem,
			prepMongoMock: func(mt *mtest.T, l *line) {
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    1,
					Message: "database general error",
					Name:    "database general error",
				}))
			},
		},
		{
			testName: "success-case",
			ctx:      context.Background(),
			id:       id,
			input: &line{
				ID:        id.GetObjectID(),
				CreatedAT: tstamp,
				Data:      "first-line",
			},
			expectedResult: &example.Line{
				ID:      id,
				Created: tstamp,
				Data:    "first-line",
			},
			prepMongoMock: func(mt *mtest.T, l *line) {
				bsonData, err := bson.Marshal(l)
				require.NoError(mt, err)

				var bsonD bson.D
				err = bson.Unmarshal(bsonData, &bsonD)
				require.NoError(mt, err)

				ns := fmt.Sprintf("%s.%s", "dbname", "lines")
				cursorResponse := mtest.CreateCursorResponse(
					1,
					ns,
					mtest.FirstBatch,
					bsonD)

				cursorEnd := mtest.CreateCursorResponse(
					0,
					ns,
					mtest.NextBatch)

				mt.AddMockResponses(cursorResponse, cursorEnd)
			},
		},
	}

	for _, c := range testCases {
		opts := mtest.NewOptions().ClientType(mtest.Mock).ShareClient(true)
		mt := mtest.New(t, opts)
		defer mt.Close()

		testName := c.testName
		ctx := c.ctx
		id := c.id
		expectedResult := c.expectedResult
		expectedError := c.expectedError

		if c.prepMongoMock != nil {
			c.prepMongoMock(mt, c.input)
		}

		mt.Run(testName, func(mt *mtest.T) {
			st := store{
				ctx:        ctx,
				collection: mt.Coll,
			}

			result, err := st.Read(ctx, id)
			assert.Equal(t, expectedResult, result)
			assert.ErrorIs(t, err, expectedError)
		})
	}
}

type configReaderMock struct {
	f func(node string) (io.Reader, error)
}

func (crm configReaderMock) Find(node string) (io.Reader, error) {
	return crm.f(node)
}

type readerMock struct {
	err error
}

func (rm readerMock) Read(p []byte) (n int, err error) {
	return 0, rm.err
}

func Test_ReadConfig(t *testing.T) {
	testCases := []struct {
		testName             string
		buildConfigReader    func(node string) (io.Reader, error)
		expectedDSN          string
		expectedDatabaseName string
		expectedCollection   string
		expectedError        error
	}{
		{
			testName: "config-read-error-case",
			buildConfigReader: func(node string) (io.Reader, error) {
				return nil, errors.New("reader error")
			},
			expectedError: ErrReadConfig,
		},
		{
			testName: "config-reader-error-case",
			buildConfigReader: func(node string) (io.Reader, error) {
				return readerMock{
					err: errors.New("read error"),
				}, nil
			},
			expectedError: ErrReadConfig,
		},
		{
			testName: "config-unmarshal-error-case",
			buildConfigReader: func(node string) (io.Reader, error) {
				return strings.NewReader("{"), nil
			},
			expectedError: ErrReadConfig,
		},
		{
			testName: "success-case",
			buildConfigReader: func(node string) (io.Reader, error) {
				return strings.NewReader(`{
					"dsn": "mongodb-dsn",
					"database": "database-name",
					"collection": "collection-name"}
				`), nil
			},
			expectedDSN:          "mongodb-dsn",
			expectedDatabaseName: "database-name",
			expectedCollection:   "collection-name",
			expectedError:        nil,
		},
	}

	for _, c := range testCases {
		name := c.testName
		expectedDSN := c.expectedDSN
		expectedDBName := c.expectedDatabaseName
		expectedCollection := c.expectedCollection
		expectedError := c.expectedError
		readerMock := configReaderMock{
			c.buildConfigReader,
		}

		t.Run(name, func(t *testing.T) {
			cnf, err := ReadConfig(readerMock)

			if cnf == nil {
				assert.Nil(t, cnf)
				assert.ErrorIs(t, err, expectedError)
			} else {
				assert.Equal(t, expectedDSN, cnf.DSN())
				assert.Equal(t, expectedDBName, cnf.Database())
				assert.Equal(t, expectedCollection, cnf.Collection())
				assert.Equal(t, expectedError, err)
			}
		})
	}
}
