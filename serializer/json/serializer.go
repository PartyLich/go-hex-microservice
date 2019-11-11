// Package json provides a concrete implementation of the RedirectSerializer interface
package json

import (
	"encoding/json"

	"github.com/PartyLich/hex-microservice/shortUrl"
	"github.com/pkg/errors"
)

type Redirect struct{}

// Encode converts a Redirect to json
func (r *Redirect) Encode(input *shortUrl.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}

// Decode converts a json []byte to a Redirect
func (r *Redirect) Decode(input []byte) (*shortUrl.Redirect, error) {
	redirect := &shortUrl.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}
