package shortUrl

// Define the repository interface our ports/adapters will need to implement

// RedirectRepository interface provides a Find method for looking up URLs based
// on their short code and a Store method for saving Redirect objects
type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
