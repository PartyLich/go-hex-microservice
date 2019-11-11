package mongo

import (
	"fmt"
	"strconv"

	"github.com/PartyLich/hex-microservice/shortUrl"
	"github.com/pkg/errors"

	"github.com/go-redis/redis"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	// create client and check connection
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, err
}

// NewRedisRepository creates an instance of RedirectRepository backed by Redis
func NewRedisRepository(redisURL string) (shortUrl.RedirectRepository, error) {
	repo := &redisRepository{}

	//generate new db client
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}

	repo.client = client
	return repo, nil
}

// utility function
func (repo *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

// Find method for looking up URLs based on their short code
func (repo *redisRepository) Find(code string) (*shortUrl.Redirect, error) {
	redirect := &shortUrl.Redirect{}
	key := repo.generateKey(code)

	// search database for code
	data, err := repo.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(shortUrl.ErrRedirectNotFound, "repository.Redirect.Find")
	}

	createdAt, err := strconv.ParseInt(data["createdAt"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

// Store method for saving Redirect objects
func (repo *redisRepository) Store(redirect *shortUrl.Redirect) error {
	key := repo.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":      redirect.Code,
		"url":       redirect.URL,
		"createdAt": redirect.CreatedAt,
	}

	// save redirect to the redis store
	_, err := repo.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.redirect.Store")
	}

	return nil
}
