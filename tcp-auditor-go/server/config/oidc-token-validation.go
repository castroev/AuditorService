package config

// OidcConfiguration contains values required for token validation middleware
type OidcConfiguration struct {
	Authority    string
	WellKnownURL string
	Scope        string
}
