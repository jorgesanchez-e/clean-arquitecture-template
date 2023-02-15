package memory

import (
	"context"
	"testing"
	"time"

	"clean-arquitecture-template/internal/domain/register"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	st := New(context.Background())

	st.stop()

	assert.NotNil(t, st.request)
	assert.NotNil(t, st.data)
	assert.NotNil(t, st.cancel)
	assert.Equal(t, 1, st.timeoutSeconds)
}

func Test_RequestType(t *testing.T) {
	testCases := []struct {
		name         string
		rtype        requestType
		expectedName string
	}{
		{
			name:         "read-request-case",
			rtype:        readRequest,
			expectedName: "read",
		},
		{
			name:         "write-request-case",
			rtype:        writeRequest,
			expectedName: "write",
		},
		{
			name:         "count-request-case",
			rtype:        countRequest,
			expectedName: "count",
		},
	}

	for _, c := range testCases {
		rtype := c.rtype
		expectedName := c.expectedName

		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, expectedName, rtype.String())
		})
	}
}

func Test_Write(t *testing.T) {
	tstamp := time.Now()

	testCases := []struct {
		name           string
		dbtimeOut      int
		input          []register.Line
		expectedOutput map[id]line
		expectedError  error
	}{
		{
			name:      "store-one-element",
			dbtimeOut: 1,
			input: []register.Line{
				{
					ID:      id("one"),
					Created: tstamp,
					Data:    "first-line",
				},
			},
			expectedOutput: map[id]line{
				"one": {
					createdAT: tstamp.Format(timeLayout),
					data:      "first-line",
				},
			},
			expectedError: nil,
		},
		{
			name:      "store-many-elements",
			dbtimeOut: 1,
			input: []register.Line{
				{
					ID:      id("one"),
					Created: tstamp,
					Data:    "first-line",
				},
				{
					ID:      id("two"),
					Created: tstamp,
					Data:    "second-line",
				},
				{
					ID:      id("three"),
					Created: tstamp,
					Data:    "third-line",
				},
			},
			expectedOutput: map[id]line{
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
			dbtimeOut: 0,
			input: []register.Line{
				{
					ID:      id("one"),
					Created: tstamp,
					Data:    "first-line",
				},
			},
			expectedError:  ErrTimeOut,
			expectedOutput: map[id]line{},
		},
	}

	for _, c := range testCases {
		ctx, cancel := context.WithCancel(context.Background())
		st := store{
			ctx:            ctx,
			cancel:         cancel,
			data:           make(map[id]line),
			request:        make(chan request),
			timeoutSeconds: c.dbtimeOut,
		}

		input := c.input
		expectedError := c.expectedError
		expectedResult := c.expectedOutput

		t.Run(c.name, func(t *testing.T) {
			st.start()

			var err error
			for _, item := range input {
				if err = st.Write(item); err != nil {
					break
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
		dbtimeOut      int
		registers      map[id]line
		searchedid     id
		expectedResult *register.Line
	}{
		{
			name:      "found-test-case",
			dbtimeOut: 1,
			registers: map[id]line{
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
			searchedid: id("two"),
			expectedResult: &register.Line{
				ID:      id("two"),
				Created: tstamp,
				Data:    "second-line",
			},
		},
		{
			name:      "not-found-test-case",
			dbtimeOut: 1,
			registers: map[id]line{
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
			searchedid:     id("x"),
			expectedResult: nil,
		},
		{
			name:           "not-found-test-case-2",
			dbtimeOut:      1,
			registers:      nil,
			searchedid:     id("x"),
			expectedResult: nil,
		},
	}

	for _, c := range testCases {
		ctx, cancel := context.WithCancel(context.Background())
		st := store{
			ctx:            ctx,
			cancel:         cancel,
			data:           c.registers,
			request:        make(chan request),
			timeoutSeconds: c.dbtimeOut,
		}

		searchedID := c.searchedid
		expectedResult := c.expectedResult

		t.Run(c.name, func(t *testing.T) {
			st.start()

			result := st.Read(searchedID)

			st.stop()

			assert.Equal(t, expectedResult, result)
		})
	}
}
