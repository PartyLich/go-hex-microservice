// Package msgpack provides a concrete implementation of the RedirectSerializer interface
package msgpack

import (
	"github.com/PartyLich/hex-microservice/shortUrl"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

// Encode converts a Redirect to msgpack
func (r *Redirect) Encode(input *shortUrl.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}

// Decode converts a msgpack []byte to a Redirect
func (r *Redirect) Decode(input []byte) (*shortUrl.Redirect, error) {
	redirect := &shortUrl.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}
