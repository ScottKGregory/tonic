package models

// Options allows for configuring tonic.
type Options struct {
	PageHeader          string             `config:"🍸, Replace the default cocktail emoji on the built in static pages"`
	DisableHomepage     bool               `config:"false, Disable the default root page"`
	DisableErrorPages   bool               `config:"false, Disable the default error pages"`
	DisableHealthProbes bool               `config:"false, Disable the default health probes"`
	Auth                AuthOptions        `config:""`
	Log                 LogOptions         `config:""`
	Backend             BackendOptions     `config:""`
	Permissions         PermissionsOptions `config:""`
}

type LogOptions struct {
	JSON         bool   `config:"false, Whether to log JSON, true, j"`
	Tag          string `config:"tonic, The log tag to use, true, t"`
	IgnoreRoutes []string
}

type AuthOptions struct {
	Disabled bool          `config:"false, Disabled the default auth system"`
	JWT      JWTOptions    `config:""`
	OIDC     OIDCOptions   `config:""`
	Cookie   CookieOptions `config:""`
}

type PermissionsOptions struct {
	Custom  []string `config:", Custom permissions to register"`
	Default []string `config:", Default permissions for new users"`
}

type JWTOptions struct {
	PrivateKey string `config:", The private key to use"`
	PublicKey  string `config:", The public key to use"`
	Duration   int64  `config:"1440, JWT token duration in minutes"`
	Audience   string `config:"tonic users, The audience to use in the token"`
	Issuer     string `config:"tonic server, The issuer to use in the token"`
}

type OIDCOptions struct {
	ClientID     string `config:", The client ID to use"`
	ClientSecret string `config:", The client secret to use"`
	Endpoint     string `config:", The endpoint to use"`
	RedirectURL  string `config:", The redirecturl to use"`
}

type CookieOptions struct {
	Name     string `config:"tonic, The name for auth cookies"`
	Path     string `config:"/, Cookie path"`
	Domain   string `config:", Cookie domain"`
	Secure   bool   `config:"true, Secure cookie"`
	HttpOnly bool   `config:"true, HTTP only"`
}

type BackendOptions struct {
	ConnectionString string `config:"mongodb://127.0.0.1:27017, The backends connection string"`
	UserCollection   string `config:"users, The backends user collection"`
	Database         string `config:"tonic, The backends database to use"`
	InMemory         bool   `config:"false, Enable to use an in memory database"`
}
