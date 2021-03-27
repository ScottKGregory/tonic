package backends

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/scottkgregory/tonic/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (mdb Mongo) connect(log *zerolog.Logger) (*mongo.Client, context.Context, context.CancelFunc) {
	client, err := mongo.NewClient(mongoOptions.Client().ApplyURI(mdb.options.ConnectionString))
	if err != nil {
		log.Error().Err(err).Msg("Error creating mongo client")
		return nil, nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		cancel()
		log.Error().Err(err).Msg("Error connecting to mongo")
		return nil, nil, nil
	}

	return client, ctx, cancel
}

func (mdb Mongo) GetUserByID(log *zerolog.Logger, id string) (user *models.User, err error) {
	user = &models.User{}
	client, ctx, cancel := mdb.connect(log)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	c := client.Database(mdb.options.Database).Collection(mdb.options.UserCollection)
	err = c.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, nil
	} else if err != nil {
		return user, err
	}

	return user, err
}

func (mdb Mongo) GetUserByOIDCSubject(log *zerolog.Logger, sub string) (user *models.User, err error) {
	user = &models.User{}
	client, ctx, cancel := mdb.connect(log)
	defer cancel()

	c := client.Database(mdb.options.Database).Collection(mdb.options.UserCollection)
	err = c.FindOne(ctx, bson.M{"claims.subject": sub}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, nil
	}

	return user, err
}

func (mdb Mongo) ListUsers(log *zerolog.Logger) (users *[]models.User, err error) {
	users = &[]models.User{}
	client, ctx, cancel := mdb.connect(log)
	defer cancel()

	c := client.Database(mdb.options.Database).Collection(mdb.options.UserCollection)
	curs, err := c.Find(ctx, bson.M{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return users, nil
	} else if err != nil {
		return users, err
	}

	err = curs.All(ctx, users)
	return users, err
}

func (mdb Mongo) SaveUser(log *zerolog.Logger, user *models.User) (err error) {
	client, ctx, cancel := mdb.connect(log)
	defer cancel()

	c := client.Database(mdb.options.Database).Collection(mdb.options.UserCollection)

	existing := &models.User{}
	u := c.FindOne(ctx, bson.M{"claims.subject": user.Claims.Subject})
	err = u.Decode(existing)
	if err == nil {
		user.ID = existing.ID
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		user.ID = primitive.NewObjectID()
	} else {
		return err
	}

	upd := bson.M{"$set": bson.M{"claims": user.Claims, "permissions": user.Permissions}}
	_, err = c.UpdateByID(ctx, user.ID, upd, mongoOptions.Update().SetUpsert(true))
	return err
}

func (mdb Mongo) Ping(log *zerolog.Logger) error {
	client, ctx, cancel := mdb.connect(log)
	defer cancel()

	err := client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return err
}
