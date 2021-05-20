package backends

import (
	"context"
	"errors"
	"time"

	"github.com/scottkgregory/tonic/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Mongo struct {
	options *models.Backend
}

var _ Backend = Mongo{}

func NewMongoBackend(options *models.Backend) *Mongo {
	return &Mongo{options}
}

func (m Mongo) connect() (*mongo.Client, context.Context, context.CancelFunc) {
	client, err := mongo.NewClient(mongoOptions.Client().ApplyURI(m.options.ConnectionString))
	if err != nil {
		return nil, nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		cancel()
		return nil, nil, nil
	}

	return client, ctx, cancel
}

func (m Mongo) CreateUser(in *models.User) (out *models.User, err error) {
	client, ctx, cancel := m.connect()
	defer cancel()

	c := client.Database(m.options.Database).Collection(m.options.UserCollection)
	_, err = c.InsertOne(ctx, in)
	return in, err
}

func (m Mongo) UpdateUser(in *models.User) (out *models.User, err error) {
	client, ctx, cancel := m.connect()
	defer cancel()

	c := client.Database(m.options.Database).Collection(m.options.UserCollection)
	upd := bson.M{"$set": bson.M{"claims": in.Claims, "permissions": in.Permissions, "deleted": in.Deleted}}
	res, err := c.UpdateOne(ctx, bson.M{"claims.subject": in.Claims.Subject}, upd)
	if res.MatchedCount == 0 {
		return nil, err
	}

	return in, err
}

func (m Mongo) GetUser(sub string) (out *models.User, err error) {
	out = &models.User{}
	client, ctx, cancel := m.connect()
	defer cancel()

	c := client.Database(m.options.Database).Collection(m.options.UserCollection)
	err = c.FindOne(ctx, bson.M{"claims.subject": sub}).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return out, err
}

func (m Mongo) ListUsers() (out []*models.User, err error) {
	out = []*models.User{}
	client, ctx, cancel := m.connect()
	defer cancel()

	c := client.Database(m.options.Database).Collection(m.options.UserCollection)
	curs, err := c.Find(ctx, bson.M{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return out, nil
	} else if err != nil {
		return out, err
	}

	err = curs.All(ctx, out)
	return out, err
}

func (m Mongo) Ping() error {
	client, ctx, cancel := m.connect()
	defer cancel()

	err := client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return err
}
