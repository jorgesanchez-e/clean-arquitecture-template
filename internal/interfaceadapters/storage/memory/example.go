package memory

import (
	"clean-arquitecture-template/internal/domain/example"
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	writeRequest requestType = iota
	readRequest
	countRequest

	timeLayout string = "2006-01-02 15:04:05"

	ErrTimeOut Err = "data store timeout"
)

var (
	storage     Store
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
	id          identifier
	input       line
	output      chan *example.Line
	count       chan *int64
}

type identifier string

func NewID() example.Identifier {
	return identifier(uuid.New().String())
}

func (id identifier) String() string {
	return string(id)
}

type line struct {
	createdAT string
	data      string
}

type Store struct {
	ctx            context.Context
	cancel         context.CancelFunc
	data           map[identifier]line
	request        chan request
	timeoutSeconds int
}

func NewExampleRepo(ctx context.Context) Store {
	storageOnce.Do(func() {
		dbCtx, dbCancel := context.WithCancel(ctx)
		storage.ctx = dbCtx
		storage.cancel = dbCancel
		storage.data = make(map[identifier]line)
		storage.request = make(chan request)
		storage.timeoutSeconds = 1

		storage.start()
	})

	return storage
}

func (s Store) start() {
	go s.run()
}

func (s Store) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case req := <-s.request:
			switch req.requestType {
			case writeRequest:
				s.data[req.id] = req.input
			case readRequest:
				req.output <- findLine(s.data, req.id)
			case countRequest:
				req.count <- count(s.data)
			}
		}
	}
}

func (s Store) stop() {
	s.cancel()
}

func (s Store) Write(ctx context.Context, n example.Line) error {
	var cancel context.CancelFunc

	if ctx == nil {
		ctx, cancel = context.WithTimeout(s.ctx, time.Duration(s.timeoutSeconds)*time.Second)
		defer cancel()
	}

	return s.write(ctx, n)
}

func (s Store) write(ctx context.Context, input example.Line) error {
	req := request{
		requestType: writeRequest,
		id:          identifier(input.ID.String()),
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

func (s Store) Read(ctx context.Context, id example.Identifier) (*example.Line, error) {
	var cancel context.CancelFunc

	if ctx == nil {
		ctx, cancel = context.WithTimeout(s.ctx, time.Duration(s.timeoutSeconds)*time.Second)
		defer cancel()
	}

	return s.read(ctx, identifier(id.String()))
}

func (s Store) read(ctx context.Context, id identifier) (*example.Line, error) {
	req := request{
		requestType: readRequest,
		id:          id,
		output:      make(chan *example.Line),
	}

	select {
	case <-ctx.Done():
		return nil, ErrTimeOut
	case s.request <- req:
		return <-req.output, nil
	}
}

func findLine(data map[identifier]line, itemID identifier) *example.Line {
	if item, exists := data[itemID]; exists {
		createdAT, _ := time.Parse(timeLayout, item.createdAT)
		return &example.Line{
			ID:      itemID,
			Created: createdAT,
			Data:    item.data,
		}
	}
	return nil
}

func (s Store) count() *int64 {
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

func count(data map[identifier]line) *int64 {
	count := new(int64)
	*count = int64(len(data))

	return count
}
