package memory

import (
	"clean-arquitecture-template/internal/domain/register"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Error(t *testing.T) {
	err := ErrTimeOut

	assert.Equal(t, "data store timeout", err.Error())
}

func Test_NewID(t *testing.T) {
	id := NewID()

	uuid, isIdentifier := id.(identifier)

	assert.True(t, isIdentifier)
	assert.NotEmpty(t, id.String())
	assert.NotEmpty(t, uuid)
}

func Test_New(t *testing.T) {
	st := New(context.Background())

	st.stop()

	assert.NotNil(t, st.request)
	assert.NotNil(t, st.data)
	assert.NotNil(t, st.cancel)
	assert.Equal(t, 1, st.timeoutSeconds)
}

func Test_RequestType(t *testing.T) {
	testCase := []struct {
		name         string
		rtype        requestType
		expectedName string
	}{
		{
			name:         "write-requesttype-case",
			rtype:        writeRequest,
			expectedName: "write",
		},
		{
			name:         "read-request-type-case",
			rtype:        readRequest,
			expectedName: "read",
		},
		{
			name:         "count-requesttype-case",
			rtype:        countRequest,
			expectedName: "count",
		},
	}

	for _, c := range testCase {
		nameCase := c.name
		expectedName := c.expectedName
		rtype := c.rtype

		t.Run(nameCase, func(t *testing.T) {
			assert.Equal(t, expectedName, rtype.String())
		})
	}
}

func Test_Write(t *testing.T) {
	tstamp := time.Now()

	testCases := []struct {
		name           string
		ctx            context.Context
		dbtimeOut      int
		input          []register.Line
		expectedOutput map[identifier]line
		expectedError  error
	}{

		{
			name:      "store-one-element",
			ctx:       context.Background(),
			dbtimeOut: 1,
			input: []register.Line{
				{
					ID:      identifier("one"),
					Created: tstamp,
					Data:    "first-line",
				},
			},
			expectedOutput: map[identifier]line{
				"one": {
					createdAT: tstamp.Format(timeLayout),
					data:      "first-line",
				},
			},
			expectedError: nil,
		},
		{
			name:      "store-one-element-nilctx",
			dbtimeOut: 1,
			ctx:       nil,
			input: []register.Line{
				{
					ID:      identifier("one"),
					Created: tstamp,
					Data:    "first-line",
				},
			},
			expectedOutput: map[identifier]line{
				"one": {
					createdAT: tstamp.Format(timeLayout),
					data:      "first-line",
				},
			},
			expectedError: nil,
		},
		{
			name:      "store-many-elements",
			ctx:       context.Background(),
			dbtimeOut: 1,
			input: []register.Line{
				{
					ID:      identifier("one"),
					Created: tstamp,
					Data:    "first-line",
				},
				{
					ID:      identifier("two"),
					Created: tstamp,
					Data:    "second-line",
				},
				{
					ID:      identifier("three"),
					Created: tstamp,
					Data:    "third-line",
				},
			},
			expectedOutput: map[identifier]line{
				"one": {
					createdAT: tstamp.Format(timeLayout),
					data:      "first-line",
				},
				"two": {
					createdAT: tstamp.Format(timeLayout),
					data:      "second-line",
				},
				"three": {
					createdAT: tstamp.Format(timeLayout),
					data:      "third-line",
				},
			},
			expectedError: nil,
		},
		{
			name:      "error",
			ctx:       context.Background(),
			dbtimeOut: 0,
			input: []register.Line{
				{
					ID:      identifier("one"),
					Created: tstamp,
					Data:    "first-line",
				},
			},
			expectedError:  ErrTimeOut,
			expectedOutput: map[identifier]line{},
		},
	}

	for _, c := range testCases {
		var storageCtx context.Context

		var cancel context.CancelFunc
		if c.ctx == nil {
			storageCtx, cancel = context.WithCancel(context.Background())
		} else {
			storageCtx, cancel = context.WithCancel(c.ctx)
		}

		st := store{
			ctx:            storageCtx,
			cancel:         cancel,
			data:           make(map[identifier]line),
			request:        make(chan request),
			timeoutSeconds: c.dbtimeOut,
		}

		input := c.input
		expectedError := c.expectedError
		expectedResult := c.expectedOutput
		timeInSec := c.dbtimeOut
		ctx := c.ctx

		t.Run(c.name, func(t *testing.T) {
			st.start()

			var err error
			for _, item := range input {
				var cancel context.CancelFunc
				var wctx context.Context

				t.Log(item)

				if ctx != nil {
					wctx, cancel = context.WithTimeout(ctx, time.Duration(timeInSec)*time.Second)
				} else {
					wctx = nil
				}

				if err = st.Write(wctx, item); err != nil {
					if wctx != nil {
						cancel()
					}
					break
				}
				if wctx != nil {
					cancel()
				}
			}

			if err == nil {
				for {
					count := st.count()
					if count != nil && *count == int64(len(input)) {
						break
					}
				}
			}

			st.stop()

			assert.Equal(t, expectedResult, st.data)
			assert.Equal(t, expectedError, err)
		})
	}
}

func Test_Read(t *testing.T) {
	tstamp := time.Date(2018, time.September, 16, 12, 0, 0, 0, time.FixedZone("", 2*60*60)).UTC()

	testCases := []struct {
		name           string
		ctx            context.Context
		dbtimeOut      int
		registers      map[identifier]line
		searchedid     identifier
		expectedResult *register.Line
		expectedError  error
	}{
		{
			name:      "found-test-case",
			ctx:       context.Background(),
			dbtimeOut: 1,
			registers: map[identifier]line{
				"one": {
					createdAT: tstamp.Format(timeLayout),
					data:      "first-line",
				},
				"two": {
					createdAT: tstamp.Format(timeLayout),
					data:      "second-line",
				},
				"three": {
					createdAT: tstamp.Format(timeLayout),
					data:      "third-line",
				},
			},
			searchedid: identifier("two"),
			expectedResult: &register.Line{
				ID:      identifier("two"),
				Created: tstamp,
				Data:    "second-line",
			},
		},
		{
			name:      "not-found-test-case",
			ctx:       context.Background(),
			dbtimeOut: 1,
			registers: map[identifier]line{
				"one": {
					createdAT: tstamp.Format(timeLayout),
					data:      "first-line",
				},
				"two": {
					createdAT: tstamp.Format(timeLayout),
					data:      "second-line",
				},
				"three": {
					createdAT: tstamp.Format(timeLayout),
					data:      "third-line",
				},
			},
			searchedid:     identifier("x"),
			expectedResult: nil,
		},
		{
			name:           "not-found-test-case-2",
			ctx:            nil,
			dbtimeOut:      1,
			registers:      nil,
			searchedid:     identifier("x"),
			expectedResult: nil,
		},
		{
			name:           "not-found-test-case-2",
			ctx:            nil,
			dbtimeOut:      0,
			registers:      nil,
			searchedid:     identifier("x"),
			expectedResult: nil,
		},
	}

	for _, c := range testCases {
		storeCtx, cancel := context.WithCancel(context.Background())

		st := store{
			ctx:            storeCtx,
			cancel:         cancel,
			data:           c.registers,
			request:        make(chan request),
			timeoutSeconds: c.dbtimeOut,
		}

		searchedID := c.searchedid
		expectedResult := c.expectedResult
		expectedError := c.expectedError
		ctx := c.ctx

		t.Run(c.name, func(t *testing.T) {
			st.start()

			result, err := st.Read(ctx, searchedID)

			st.stop()

			assert.Equal(t, expectedResult, result)
			assert.Equal(t, expectedError, err)
		})
	}
}
