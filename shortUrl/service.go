package shortUrl

// RedirectService interface provides a Find method for looking up URLs based
// on their short code and a Store method for saving Redirect objects
type RedirectService interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
