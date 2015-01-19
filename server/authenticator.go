package server

type Authenticator interface {
	Authenticate(string, string) bool
}

type AuthenticatorFunc func(string, string) bool

type simpleAuthenticator struct {
	f AuthenticatorFunc
}

func NewSimpleAuthenticator(f AuthenticatorFunc) Authenticator {
	return &simpleAuthenticator{
		f: f,
	}
}

func (sa *simpleAuthenticator) Authenticate(username string, password string) bool {
	return sa.f(username, password)
}
