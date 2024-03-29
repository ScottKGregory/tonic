package backends

import (
	"context"
	"errors"

	"github.com/scottkgregory/tonic/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type Mongo struct {
	config *models.BackendConfig
	client *mongo.Client
}

var _ Backend = Mongo{}

func NewMongoBackend(ctx context.Context, config *models.BackendConfig) (*Mongo, error) {
	client, err := mongo.NewClient(mongoOptions.Client().ApplyURI(config.ConnectionString))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &Mongo{config, client}, nil
}

func (m Mongo) CreateUser(ctx context.Context, in *models.User) (out *models.User, err error) {
	c := m.client.Database(m.config.Database).Collection(m.config.UserCollection)
	_, err = c.InsertOne(ctx, in)
	return in, err
}

func (m Mongo) UpdateUser(ctx context.Context, in *models.User) (out *models.User, err error) {
	c := m.client.Database(m.config.Database).Collection(m.config.UserCollection)
	upd := bson.M{"$set": bson.M{"claims": in.Claims, "permissions": in.Permissions, "deleted": in.Deleted}}
	res, err := c.UpdateOne(ctx, bson.M{"claims.subject": in.Claims.Subject}, upd)
	if res.MatchedCount == 0 {
		return nil, err
	}

	return in, err
}

func (m Mongo) GetUser(ctx context.Context, sub string) (out *models.User, err error) {
	out = &models.User{}
	c := m.client.Database(m.config.Database).Collection(m.config.UserCollection)
	err = c.FindOne(ctx, bson.M{"claims.subject": sub}).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	return out, err
}

func (m Mongo) ListUsers(ctx context.Context) (out []*models.User, err error) {
	out = []*models.User{}
	c := m.client.Database(m.config.Database).Collection(m.config.UserCollection)
	curs, err := c.Find(ctx, bson.M{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return out, nil
	} else if err != nil {
		return out, err
	}

	err = curs.All(ctx, &out)
	return out, err
}

func (m Mongo) Ping(ctx context.Context) error {
	err := m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return err
}
