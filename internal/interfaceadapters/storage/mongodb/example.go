package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"clean-arquitecture-template/internal/domain/example"
)

const (
	ErrIdentifyer   mongoError = "invalid mongodb identifyer"
	ErrDataInserted mongoError = "db error on insert-one"
	ErrMongoSystem  mongoError = "database error"

	objectIDRegexpFormat string = `ObjectID\("([a-zA-Z0-9]+)"\)`
)

type mongoError string

func (me mongoError) Error() string {
	return string(me)
}

type Config interface {
	GetDSN() string
	DatabaseName() string
	TableName() string
}

type Identifier primitive.ObjectID

func NewIdentifier() Identifier {
	return Identifier{}
}

func (id Identifier) NewID() example.Identifier {
	return Identifier(primitive.NewObjectID())
}

func (id Identifier) String() string {
	rx := regexp.MustCompile(objectIDRegexpFormat)
	r := rx.FindStringSubmatch(primitive.ObjectID(id).String())

	if len(r) != 2 {
		return ""
	}

	return r[1]
}

func (id Identifier) GetObjectID() primitive.ObjectID {
	return primitive.ObjectID(id)
}

func (id Identifier) ParseID(key string) (example.Identifier, error) {
	oid, err := primitive.ObjectIDFromHex(key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrIdentifyer)
	}

	return Identifier(oid), nil
}

type mongoCollection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
}

type store struct {
	ctx        context.Context
	collection mongoCollection
}

func NewExampleRepo(ctx context.Context, conf Config) store {
	clientOptions := options.Client().ApplyURI(conf.GetDSN())
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	return store{
		ctx:        ctx,
		collection: client.Database(conf.DatabaseName()).Collection(conf.TableName()),
	}
}

type line struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAT time.Time          `bson:"created_at"`
	Data      string             `bson:"data"`
}

func newLine(id primitive.ObjectID, createdAT time.Time, data string) line {
	return line{
		ID:        id,
		CreatedAT: createdAT,
		Data:      data,
	}
}

func (l *line) registerLine() *example.Line {
	if l == nil {
		return nil
	}

	return &example.Line{
		ID:      Identifier(l.ID),
		Created: l.CreatedAT,
		Data:    l.Data,
	}
}

func (s store) Write(ctx context.Context, wline example.Line) error {
	if ctx == nil {
		ctx = s.ctx
	}

	if id, is := wline.ID.(Identifier); !is {
		return ErrIdentifyer
	} else {
		return s.write(ctx, newLine(id.GetObjectID(), wline.Created, wline.Data))
	}
}

func (s store) write(ctx context.Context, nline line) error {
	_, err := s.collection.InsertOne(s.ctx, nline)
	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), ErrDataInserted)
	}

	return err
}

func (s store) Read(ctx context.Context, id example.Identifier) (*example.Line, error) {
	if ctx == nil {
		ctx = s.ctx
	}

	if id, is := id.(Identifier); !is {
		return nil, ErrIdentifyer
	} else {
		r, err := s.read(ctx, id.GetObjectID())
		return r, err
	}
}

func (s store) read(ctx context.Context, id primitive.ObjectID) (*example.Line, error) {
	payload := new(line)
	filter := bson.D{{Key: "_id", Value: id}}

	if err := s.collection.FindOne(ctx, filter).Decode(payload); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", err.Error(), ErrMongoSystem)
	}

	return payload.registerLine(), nil
}
