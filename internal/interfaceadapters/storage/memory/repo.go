package memory

import (
	"context"
	"sync"
	"time"

	"clean-arquitecture-template/internal/domain/register"
)

const (
	writeRequest requestType = iota
	readRequest
	countRequest

	timeLayout string = "2006-01-02 15:04:05"

	ErrTimeOut Err = "data store timeout"
)

var (
	storage     store
	storageOnce sync.Once
)

type Err string

func (e Err) Error() string {
	return string(e)
}

type requestType int

func (rt requestType) String() string {
	return []string{"write", "read", "count"}[rt]
}

type request struct {
	requestType requestType
	id          string
	input       line
	output      chan *register.Line
	count       chan *int64
}

type id string

func (id id) String() string {
	return string(id)
}

type line struct {
	createdAT string
	data      string
}

type store struct {
	ctx            context.Context
	cancel         context.CancelFunc
	data           map[id]line
	request        chan request
	timeoutSeconds int
}

func New(ctx context.Context) store {
	storageOnce.Do(func() {
		dbCtx, dbCancel := context.WithCancel(ctx)
		storage.ctx = dbCtx
		storage.cancel = dbCancel
		storage.data = make(map[id]line)
		storage.request = make(chan request)
		storage.timeoutSeconds = 1

		storage.start()
	})

	return storage
}

func (s store) start() {
	go s.run()
}

func (s store) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case req := <-s.request:
			switch req.requestType {
			case writeRequest:
				s.data[id(req.id)] = req.input
			case readRequest:
				req.output <- findLine(s.data, req.id)
			case countRequest:
				req.count <- count(s.data)
			}
		}
	}
}

func (s store) stop() {
	s.cancel()
}

func (s store) Write(n register.Line) error {
	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(s.timeoutSeconds)*time.Second)
	defer cancel()

	return s.write(ctx, n)
}

func (s store) write(ctx context.Context, input register.Line) error {
	req := request{
		requestType: writeRequest,
		id:          input.ID.String(),
		input: line{
			createdAT: input.Created.Format(timeLayout),
			data:      input.Data,
		},
	}

	if err := ctx.Err(); err != nil {
		return ErrTimeOut
	}

	select {
	case <-ctx.Done():
		return ErrTimeOut
	case s.request <- req:
		return nil
	}
}

func (s store) Read(id register.Identifier) *register.Line {
	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(s.timeoutSeconds)*time.Second)
	defer cancel()

	return s.read(ctx, id.String())
}

func (s store) read(ctx context.Context, id string) *register.Line {
	req := request{
		requestType: readRequest,
		id:          id,
		output:      make(chan *register.Line),
	}

	select {
	case <-ctx.Done():
		return nil
	case s.request <- req:
		return <-req.output
	}
}

func findLine(data map[id]line, itemID string) *register.Line {
	if item, exists := data[id(itemID)]; exists {
		createdAT, _ := time.Parse(timeLayout, item.createdAT)
		return &register.Line{
			ID:      id(itemID),
			Created: createdAT,
			Data:    item.data,
		}
	}
	return nil
}

func (s store) count() *int64 {
	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(s.timeoutSeconds))
	defer cancel()

	req := request{
		requestType: countRequest,
		count:       make(chan *int64),
	}

	select {
	case <-ctx.Done():
		return nil
	case s.request <- req:
	}

	count := <-req.count

	return count
}

func count(data map[id]line) *int64 {
	count := new(int64)
	*count = int64(len(data))

	return count
}
