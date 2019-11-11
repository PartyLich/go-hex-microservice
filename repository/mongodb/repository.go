package mongo

import (
	"context"
	"time"

	"github.com/PartyLich/hex-microservice/shortUrl"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	client   *mongo.Client
	database string // name of the database
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	// perform read to check our db access
	// err = client.Ping(ctx, readpref.Primary())
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, err
}

// NewMongoRepository creates an instance of RedirectRepository backed by mongoDB
func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (shortUrl.RedirectRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}

	//generate new db client
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepository")
	}

	repo.client = client
	return repo, nil
}

// Find method for looking up URLs based on their short code
func (repo *mongoRepository) Find(code string) (*shortUrl.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.timeout)
	defer cancel()

	redirect := &shortUrl.Redirect{}
	collection := repo.client.Database(repo.database).Collection("redirects")

	// search database for code
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortUrl.ErrRedirectNotFound, "repository.redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	return redirect, nil
}

// Store method for saving Redirect objects
func (repo *mongoRepository) Store(redirect *shortUrl.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), repo.timeout)
	defer cancel()

	collection := repo.client.Database(repo.database).Collection("redirects")

	// save redirect to the db
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":      redirect.Code,
			"url":       redirect.URL,
			"createdAt": redirect.CreatedAt,
		})
	if err != nil {
		return errors.Wrap(err, "repository.redirect.Store")
	}

	return nil
}
