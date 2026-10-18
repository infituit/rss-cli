package types

// AuthResult is the provider-neutral result of an authentication check.
type AuthResult struct {
	APIVersion          int
	Authenticated       bool
	LastRefreshedOnTime string
}
