package services

import (
	"crypto/rsa"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
	"github.com/rs/zerolog"
	"github.com/scottkgregory/tonic/pkg/backends"
	"github.com/scottkgregory/tonic/pkg/constants"
	"github.com/scottkgregory/tonic/pkg/helpers"
	"github.com/scottkgregory/tonic/pkg/models"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	errorPage    = "/error/500"
	successPage  = "/"
	BearerPrefix = "Bearer "
)

type AuthService struct {
	state       string
	config      *oauth2.Config
	provider    *oidc.Provider
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	backend     *backends.Mongo
	authOptions *models.Auth
	permService *PermissionService
}

func NewAuthService(authOptions *models.Auth, permService *PermissionService, backendOptions *models.Backend) *AuthService {
	log := helpers.GetLogger()
	var err error
	privateKey, err := helpers.ParsePrivateKey(authOptions.JWT.PrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading private key")
		return nil
	}

	publicKey, err := helpers.ParsePublicKey(authOptions.JWT.PublicKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading public key")
		return nil
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, authOptions.OIDC.Endpoint)
	if err != nil {
		log.Fatal().Err(err).Msg("Error setting up OIDC provider")
		return nil
	}

	return &AuthService{state: "Tonic", config: &oauth2.Config{
		ClientID:     authOptions.OIDC.ClientID,
		ClientSecret: authOptions.OIDC.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  authOptions.OIDC.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	},
		provider:    provider,
		privateKey:  privateKey,
		publicKey:   publicKey,
		backend:     backends.NewMongoBackend(backendOptions),
		authOptions: authOptions,
		permService: permService,
	}
}

func (s *AuthService) Login(c *gin.Context) {
	payload := []byte(s.state)

	encrypted, err := jwe.Encrypt(payload, jwa.RSA1_5, s.publicKey, jwa.A128CBC_HS256, jwa.NoCompress)
	if err != nil {
		s.retError(c, err, "Error encrypting JWE")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, s.config.AuthCodeURL(string(encrypted)))
}

func (s *AuthService) Logout(c *gin.Context) {
	c.SetCookie(
		s.authOptions.Cookie.Name,
		"",
		-1,
		s.authOptions.Cookie.Path,
		s.authOptions.Cookie.Domain,
		s.authOptions.Cookie.Secure,
		s.authOptions.Cookie.HttpOnly,
	)
	c.Redirect(http.StatusTemporaryRedirect, successPage)
}

func (s *AuthService) Callback(c *gin.Context) {
	log := helpers.GetLogger()

	state := c.Query("state")

	decrypted, err := jwe.Decrypt([]byte(state), jwa.RSA1_5, s.privateKey)
	if err != nil {
		s.retError(c, err, "Error decrypting JWE")
		return
	}

	if string(decrypted) != s.state {
		s.retError(c, err, "Decrypted state does not match expected")
		return
	}

	oauth2Token, err := s.config.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		s.retError(c, err, "Error exchanging OIDC token")
		return
	}

	userInfo, err := s.provider.UserInfo(c.Request.Context(), oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		s.retError(c, err, "Error getting user info")
		return
	}

	um, err := s.backend.GetUserByOIDCSubject(log, userInfo.Subject)
	if err != nil {
		s.retError(c, err, "Error getting user from backend")
		return
	}

	um.Claims = models.StandardClaims{
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
		s.retError(c, err, "Error parsing user claims")
		return
	}

	err = s.backend.SaveUser(log, um)
	if err != nil {
		s.retError(c, err, "Error saving user info")
		return
	}

	t := openid.New()

	if err := t.Set(jwt.IssuerKey, s.authOptions.JWT.Issuer); err != nil {
		s.retError(c, err, "Error setting %s on JWT", jwt.IssuerKey)
		return
	}

	if err := t.Set(jwt.SubjectKey, userInfo.Subject); err != nil {
		s.retError(c, err, "Error setting %s on JWT", jwt.SubjectKey)
		return
	}

	if err := t.Set(jwt.AudienceKey, s.authOptions.JWT.Audience); err != nil {
		s.retError(c, err, "Error setting %s on JWT", jwt.AudienceKey)
		return
	}

	if err := t.Set(jwt.IssuedAtKey, time.Now()); err != nil {
		s.retError(c, err, "Error setting %s on JWT", jwt.IssuedAtKey)
		return
	}

	if err := t.Set(jwt.ExpirationKey, time.Now().Add(time.Duration(s.authOptions.JWT.Duration)*time.Minute)); err != nil {
		s.retError(c, err, "Error setting %s on JWT", jwt.ExpirationKey)
		return
	}

	if err := t.Set(constants.PermissionsKey, um.Permissions); err != nil {
		s.retError(c, err, "Error setting %s on JWT", constants.PermissionsKey)
		return
	}

	signed, err := jwt.Sign(t, jwa.RS256, s.privateKey)
	if err != nil {
		s.retError(c, err, "Error signing JWT")
		return
	}

	c.SetCookie(
		s.authOptions.Cookie.Name,
		string(signed),
		int(s.authOptions.JWT.Duration)*60,
		s.authOptions.Cookie.Path,
		s.authOptions.Cookie.Domain,
		s.authOptions.Cookie.Secure,
		s.authOptions.Cookie.HttpOnly,
	)

	c.Redirect(http.StatusTemporaryRedirect, successPage)
}

func (s *AuthService) Token(c *gin.Context) {
	log := helpers.GetLogger()

	var tok string
	var err error
	if c.GetString(constants.AuthMethodKey) == constants.Cookie {
		tok, err = c.Cookie(s.authOptions.Cookie.Name)
		if err != nil {
			log.Error().Err(err).Msg("Error getting cookie")
			helpers.APIErrorResponse(c, http.StatusInternalServerError, "Error getting cookie")
			return
		}
	} else {
		tok = c.GetHeader(constants.Authorization)
		tok = strings.TrimPrefix(tok, BearerPrefix)
	}

	token, err := jwt.Parse(
		[]byte(tok),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, s.publicKey),
	)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing JWT")
		helpers.APIErrorResponse(c, http.StatusInternalServerError, "Error parsing JWT")
		return
	}

	helpers.APISuccessResponse(c, &models.Token{
		Token:  tok,
		Expiry: token.Expiration(),
	})
}

func (s *AuthService) Me(c *gin.Context) {
	log := helpers.GetLogger()

	sub := c.Keys[constants.SubjectKey].(string)
	u, err := s.backend.GetUserByOIDCSubject(helpers.GetLogger(c), sub)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user info")
		helpers.APIErrorResponse(c, http.StatusInternalServerError, "Unable to retrieve user")
		return
	}

	helpers.APISuccessResponse(c, u)
}

func (s *AuthService) retError(c *gin.Context, err error, message string, vals ...interface{}) {
	log := helpers.GetLogger(c)
	log.Error().Err(err).Msgf(message, vals...)
	c.Redirect(http.StatusTemporaryRedirect, errorPage)
}

func (s *AuthService) Renew(old string, log *zerolog.Logger) (new string, err error) {
	parsed, err := jwt.Parse(
		[]byte(old),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, s.publicKey),
	)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing JWT")
		return new, err
	}

	if err = parsed.Set(jwt.IssuedAtKey, time.Now()); err != nil {
		log.Error().Err(err).Msgf("Error setting %s on JWT", jwt.IssuedAtKey)
		return
	}

	if err = parsed.Set(jwt.ExpirationKey, time.Now().Add(time.Duration(s.authOptions.JWT.Duration)*time.Minute)); err != nil {
		log.Error().Err(err).Msgf("Error setting %s on JWT", jwt.ExpirationKey)
		return
	}

	signed, err := jwt.Sign(parsed, jwa.RS256, s.privateKey)
	if err != nil {
		log.Error().Err(err).Msgf("Error signing JWT")
		return
	}

	return string(signed), nil
}

func (s *AuthService) Verify(tok string, log *zerolog.Logger) (valid bool, subject string, expiry time.Time, permissions []string) {
	var perms []string
	token, err := jwt.Parse(
		[]byte(tok),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, s.publicKey),
	)
	if err != nil {
		return false, "", time.Now(), perms
	}

	p, b := token.Get(constants.PermissionsKey)
	if b {
		for _, x := range p.([]interface{}) {
			perms = append(perms, x.(string))
		}
	}

	return true, token.Subject(), token.Expiration(), perms
}
