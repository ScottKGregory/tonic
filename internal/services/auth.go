package services

import (
	"context"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
	"github.com/rs/zerolog"
	tonicErrors "github.com/scottkgregory/tonic/internal/api/errors"
	"github.com/scottkgregory/tonic/internal/constants"
	"github.com/scottkgregory/tonic/internal/helpers"
	"github.com/scottkgregory/tonic/internal/models"
	pkgModels "github.com/scottkgregory/tonic/pkg/models"
	"golang.org/x/oauth2"
)

// AuthService contains auth related operations
type AuthService struct {
	state       string
	log         *zerolog.Logger
	userService *UserService
	permService *PermissionsService
	options     *models.AuthOptions
	authConfig  *oauth2.Config
	provider    *oidc.Provider
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
}

// NewAuthService configures a new instance of AuthService
func NewAuthService(log *zerolog.Logger, userService *UserService, permService *PermissionsService, options *models.AuthOptions) *AuthService {
	var err error
	privateKey, err := helpers.ParsePrivateKey(options.JWT.PrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading private key")
		return nil
	}

	publicKey, err := helpers.ParsePublicKey(options.JWT.PublicKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading public key")
		return nil
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, options.OIDC.Endpoint)
	if err != nil {
		log.Fatal().Err(err).Msg("Error setting up OIDC provider")
		return nil
	}

	authConfig := &oauth2.Config{
		ClientID:     options.OIDC.ClientID,
		ClientSecret: options.OIDC.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  options.OIDC.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	state := "tonic"

	return &AuthService{
		state,
		log,
		userService,
		permService,
		options,
		authConfig,
		provider,
		privateKey,
		publicKey,
	}
}

// Login gets the OIDC login URL for the given provider
func (s *AuthService) Login(provider string) (redirect string, err error) {
	payload := []byte(s.state)

	encrypted, err := jwe.Encrypt(payload, jwa.RSA1_5, s.publicKey, jwa.A128CBC_HS256, jwa.NoCompress)
	if err != nil {
		return "", err
	}

	return s.authConfig.AuthCodeURL(string(encrypted)), err
}

// Callback processes the OIDC flow return values
func (s *AuthService) Callback(ctx context.Context, provider, state, code, callbackErr, errDescription string) (token string, err error) {
	if helpers.IsEmptyOrWhitespace(code) ||
		helpers.IsEmptyOrWhitespace(state) ||
		!helpers.IsEmptyOrWhitespace(callbackErr) ||
		!helpers.IsEmptyOrWhitespace(errDescription) {
		return "", tonicErrors.NewUnauthorisedError()
	}

	decrypted, err := jwe.Decrypt([]byte(state), jwa.RSA1_5, s.privateKey)
	if err != nil {
		return "", err
	}

	if string(decrypted) != s.state {
		return "", err
	}

	oauth2Token, err := s.authConfig.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	userInfo, err := s.provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		return "", err
	}

	um, err := s.userService.GetUser(userInfo.Subject)
	if errors.Is(err, &tonicErrors.NotFoundErr{}) {
		um, err = s.userService.CreateUser(&pkgModels.User{
			Claims: pkgModels.StandardClaims{
				Subject: userInfo.Subject,
			},
		})
	}

	if err != nil {
		return "", err
	}

	um.Claims = pkgModels.StandardClaims{
		Subject:       userInfo.Subject,
		Profile:       userInfo.Profile,
		Email:         userInfo.Email,
		EmailVerified: userInfo.EmailVerified,
	}

	if len(um.Permissions) == 0 {
		um.Permissions = s.permService.DefaultPermissions()
	}

	err = userInfo.Claims(&um.Claims)
	if err != nil {
		return "", err
	}

	um, err = s.userService.UpdateUser(um, um.Claims.Subject)
	if err != nil {
		return "", err
	}

	t, err := s.createToken(um)
	if err != nil {
		return "", err
	}

	signed, err := jwt.Sign(t, jwa.RS256, s.privateKey)
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

// Token generates an auth token for the given user
func (s *AuthService) Token(subject string) (token *pkgModels.Token, err error) {
	if helpers.IsEmptyOrWhitespace(subject) {
		return nil, tonicErrors.NewUnauthorisedError()
	}

	user, err := s.userService.GetUser(subject)
	if err != nil {
		return nil, err
	}

	oidcTok, err := s.createToken(user)
	if err != nil {
		return nil, err
	}

	signed, err := jwt.Sign(oidcTok, jwa.RS256, s.privateKey)
	if err != nil {
		return nil, err
	}

	return &pkgModels.Token{
		Token:  string(signed),
		Expiry: oidcTok.Expiration(),
	}, nil
}

// Verify parses and verifies the provided token
func (s *AuthService) Verify(tok string) (bool, jwt.Token) {
	token, err := jwt.Parse(
		[]byte(tok),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, s.publicKey),
	)
	if err != nil {
		return false, nil
	}

	return true, token
}

func (s *AuthService) createToken(user *pkgModels.User) (token openid.Token, err error) {
	t := openid.New()

	if err := t.Set(jwt.IssuerKey, s.options.JWT.Issuer); err != nil {
		return nil, err
	}

	if err := t.Set(jwt.SubjectKey, user.Claims.Subject); err != nil {
		return nil, err
	}

	if err := t.Set(jwt.AudienceKey, s.options.JWT.Audience); err != nil {
		return nil, err
	}

	if err := t.Set(jwt.IssuedAtKey, time.Now()); err != nil {
		return nil, err
	}

	exp := time.Now().Add(time.Duration(s.options.JWT.Duration) * time.Minute).UTC()
	if err := t.Set(jwt.ExpirationKey, exp); err != nil {
		return nil, err
	}

	if err := t.Set(constants.PermissionsKey, user.Permissions); err != nil {
		return nil, err
	}

	return t, err
}
