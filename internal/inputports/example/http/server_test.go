package http

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type configReaderMock struct {
	f func(node string) (io.Reader, error)
}

func (cr configReaderMock) Find(node string) (io.Reader, error) {
	return cr.f(node)
}

type readerMock struct {
	err error
}

func (rm readerMock) Read(p []byte) (n int, err error) {
	return 0, rm.err
}

func Test_ReadConfig(t *testing.T) {
	testCases := []struct {
		testName        string
		configReader    func(node string) (io.Reader, error)
		expectedAddress string
		expectedError   error
	}{
		{
			testName: "error-read-config-case",
			configReader: func(node string) (io.Reader, error) {
				return nil, errors.New("some-error")
			},
			expectedAddress: "",
			expectedError:   ErrReadConfig,
		},
		{
			testName: "reader-config-error-case",
			configReader: func(node string) (io.Reader, error) {
				r := readerMock{err: errors.New("some-error")}
				return r, nil
			},
			expectedAddress: "",
			expectedError:   ErrReadConfig,
		},
		{
			testName: "unmarshal-error-case",
			configReader: func(node string) (io.Reader, error) {
				r := strings.NewReader("{")
				return r, nil
			},
			expectedAddress: "",
			expectedError:   ErrReadConfig,
		},
		{
			testName: "success-case",
			configReader: func(node string) (io.Reader, error) {
				r := strings.NewReader(`{
					"address": "127.0.0.1",
					"port":    "8080"					
				}`)
				return r, nil
			},
			expectedAddress: "127.0.0.1:8080",
			expectedError:   nil,
		},
	}

	for _, c := range testCases {
		name := c.testName
		expectedAddres := c.expectedAddress
		expectedError := c.expectedError

		mock := configReaderMock{
			f: c.configReader,
		}

		t.Run(name, func(t *testing.T) {
			config, err := ReadConfig(mock)

			if config == nil {
				assert.ErrorIs(t, err, expectedError)
			} else {
				assert.Equal(t, expectedAddres, config.Address())
				assert.ErrorIs(t, err, expectedError)
			}
		})
	}
}
