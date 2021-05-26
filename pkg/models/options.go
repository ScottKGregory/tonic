package models

// Options allows for configuring tonic.
type Options struct {
	PageHeader          string `config:"🍸, Replace the default cocktail emoji on the built in static pages"`
	DisableHomepage     bool   `config:"false, Disable the default root page"`
	DisableErrorPages   bool   `config:"false, Disable the default error pages"`
	DisableHealthProbes bool   `config:"false, Disable the default health probes"`
	Auth                `config:""`
	Log                 `config:""`
	Backend             `config:""`
}

type Log struct {
	JSON         bool   `config:"false, Whether to log JSON, true, j"`
	Tag          string `config:"tonic, The log tag to use, true, t"`
	IgnoreRoutes []string
}

type Auth struct {
	Disabled bool   `config:"false, Disabled the default auth system"`
	JWT      JWT    `config:""`
	OIDC     OIDC   `config:""`
	Cookie   Cookie `config:""`
}

type JWT struct {
	PrivateKey string `config:", The private key to use"`
	PublicKey  string `config:", The public key to use"`
	Duration   int64  `config:"1440, JWT token duration in minutes"`
	Audience   string `config:"tonic users, The audience to use in the token"`
	Issuer     string `config:"tonic server, The issuer to use in the token"`
}

type OIDC struct {
	ClientID     string `config:", The client ID to use"`
	ClientSecret string `config:", The client secret to use"`
	Endpoint     string `config:", The endpoint to use"`
	RedirectURL  string `config:", The redirecturl to use"`
}

type Cookie struct {
	Name     string `config:"tonic, The name for auth cookies"`
	Path     string `config:"/, Cookie path"`
	Domain   string `config:", Cookie domain"`
	Secure   bool   `config:"true, Secure cookie"`
	HttpOnly bool   `config:"true, HTTP only"`
}

type Backend struct {
	ConnectionString string `config:"mongodb://127.0.0.1:27017, The backends connection string"`
	UserCollection   string `config:"users, The backends user collection"`
	Database         string `config:"tonic, The backends database to use"`
	InMemory         bool   `config:"false, Enable to use an in memory database"`
}
