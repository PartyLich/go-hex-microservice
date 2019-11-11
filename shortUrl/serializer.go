package shortUrl

// Define the serializer interface our ports/adapters will need to implement

// RedirectSerializer provides methods to serialize and deserialize Redirect objects
type RedirectSerializer interface {
	Encode(input *Redirect) ([]byte, error)
	Decode(input []byte) (*Redirect, error)
}
