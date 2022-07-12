package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os/user"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/shumin1027/otpd/pkg/pam"
)

type AuthSchema string

const (
	AuthSchemaNone   = "None"
	AuthSchemaBasic  = "Basic"
	AuthSchemaBearer = "Bearer"
	AuthSchemaOther  = "Other"
)

type AuthResult struct {
	AuthSchema AuthSchema
	Token      jwt.Token
	User       *user.User
}

// Config defines the config for BasicAuth middleware
type Config struct {
	// AuthSchemas to be used in the Authorization header.
	// Optional. Default: "Bearer, Basic".
	AuthSchemas []AuthSchema

	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool

	// SuccessHandler defines a function which is executed for a valid token.
	// Optional. Default: nil
	SuccessHandler fiber.Handler

	// ErrorHandler defines a function which is executed for an invalid token.
	// It may be used to define a custom JWT error.
	// Optional. Default: 401 Invalid or expired JWT
	ErrorHandler fiber.ErrorHandler

	NoneAuthHandler fiber.Handler

	JWTParseOptions []jwt.ParseOption

	// Signing key to validate token. Used as fallback if SigningKeys has length 0.
	// Required. This or SigningKeys.
	SigningKey jwk.Key

	// Signing method, used to check token signing method.
	// Optional. Default: "HS256".
	// Possible values: "HS256", "HS384", "HS512", "ES256", "ES384", "ES512", "RS256", "RS384", "RS512"
	SignatureAlgorithm jwa.SignatureAlgorithm

	// Context key to store user information from the token into context.
	// Optional. Default: "auth".
	ContextKey string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "param:<name>"
	// - "cookie:<name>"
	TokenLookup string
}

// New auth middleware
func New(config ...Config) fiber.Handler {
	// Init config
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	// 默认handler定义
	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}
	}
	if cfg.NoneAuthHandler == nil {
		cfg.NoneAuthHandler = func(c *fiber.Ctx) error {
			return cfg.ErrorHandler(c, errors.New("authenticate is none"))
		}
	}

	// jwt token 解析认证项
	if cfg.JWTParseOptions == nil {
		cfg.JWTParseOptions = []jwt.ParseOption{}
	}
	if cfg.SignatureAlgorithm == "" {
		cfg.SignatureAlgorithm = jwa.HS256
	}
	if cfg.TokenLookup == "" {
		cfg.TokenLookup = "header:" + fiber.HeaderAuthorization
	}

	if cfg.ContextKey == "" {
		cfg.ContextKey = "auth"
	}

	if cfg.AuthSchemas == nil {
		cfg.AuthSchemas = []AuthSchema{AuthSchemaBasic, AuthSchemaBearer}
	}

	// Initialize 初始化从不同方式获取授权信息的函数
	extractors := make([]func(c *fiber.Ctx) (string, error), 0)
	rootParts := strings.Split(cfg.TokenLookup, ",")
	for _, rootPart := range rootParts {
		parts := strings.Split(strings.TrimSpace(rootPart), ":")

		switch parts[0] {
		case "header":
			extractors = append(extractors, authFromHeader(parts[1]))
		case "query":
			extractors = append(extractors, authFromQuery(parts[1]))
		case "param":
			extractors = append(extractors, authFromParam(parts[1]))
		case "cookie":
			extractors = append(extractors, authFromCookie(parts[1]))
		}
	}

	// Return middleware handler
	return func(c *fiber.Ctx) error {
		// Filter request to skip middleware
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		// Get auth schema from request
		authSchema := getAuthSchema(c)

		if authSchema == AuthSchemaNone {
			// 执行NoneAuthHandler
			return cfg.NoneAuthHandler(c)

		}
		// 判断当前AuthSchema是否被允许
		allow := false
		for _, schema := range cfg.AuthSchemas {
			if schema == authSchema {
				allow = true
				break
			}
		}
		if !allow {
			return cfg.ErrorHandler(c, fmt.Errorf("AuthSchema:%s not allowed", authSchema))
		}

		var auth string
		var err error

		// Get auth information from request
		for _, extractor := range extractors {
			auth, err = extractor(c)
			if auth != "" && err == nil {
				break
			}
		}

		if err != nil {
			return cfg.ErrorHandler(c, err)
		}

		switch authSchema {
		case AuthSchemaBasic:
			{
				u, err := basicAuthVerifyByPam(auth)
				if err != nil {
					return cfg.ErrorHandler(c, err)
				}
				res := &AuthResult{
					AuthSchema: AuthSchemaBasic,
					User:       u,
				}
				c.Locals(cfg.ContextKey, res)
				return cfg.SuccessHandler(c)
			}
		case AuthSchemaBearer:
			{
				token, err := bearerAuthVerifyByJwt(auth, cfg.SignatureAlgorithm, cfg.SigningKey, cfg.JWTParseOptions...)
				if err != nil {
					return cfg.ErrorHandler(c, err)
				}
				res := &AuthResult{
					AuthSchema: AuthSchemaBearer,
					Token:      token,
				}
				c.Locals(cfg.ContextKey, res)
				return cfg.SuccessHandler(c)
			}
		default:
			{
				return cfg.ErrorHandler(c, fmt.Errorf("AuthSchema:%s not allowed", authSchema))
			}
		}
	}
}

// jwtFromHeader returns a function that extracts token from the request header.
func authFromHeader(header string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		authSchema := getAuthSchema(c)
		auth := c.Get(header)
		l := len(authSchema)
		if len(auth) > l+1 && strings.EqualFold(auth[:l], string(authSchema)) {
			return auth[l+1:], nil
		}
		return "", errors.New("missing or malformed auth info")
	}
}

// jwtFromQuery returns a function that extracts token from the query string.
func authFromQuery(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Query(param)
		if token == "" {
			return "", errors.New("missing or malformed auth info")
		}
		return token, nil
	}
}

// jwtFromParam returns a function that extracts token from the url param string.
func authFromParam(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Params(param)
		if token == "" {
			return "", errors.New("missing or malformed auth info")
		}
		return token, nil
	}
}

// jwtFromCookie returns a function that extracts token from the named cookie.
func authFromCookie(name string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Cookies(name)
		if token == "" {
			return "", errors.New("missing or malformed auth info")
		}
		return token, nil
	}
}

func basicAuthVerifyByPam(auth string) (u *user.User, err error) {
	username, password, ok := parseBasicAuth(auth)
	if !ok {
		return nil, errors.New("invalid username or password")
	}
	u, err = user.Lookup(username)
	if err != nil {
		return nil, fmt.Errorf("user:%s not found", username)
	}
	p, err := pam.Authenticate(username, password)
	if err != nil {
		return nil, err
	}
	if p == pam.PamSuccess {
		return u, nil
	}
	return nil, errors.New("invalid username or password")
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	c, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

func bearerAuthVerifyByJwt(auth string, alg jwa.SignatureAlgorithm, key jwk.Key, opts ...jwt.ParseOption) (token jwt.Token, err error) {
	opts = append(opts, jwt.WithValidate(true))
	opts = append(opts, jwt.WithVerify(alg, key))
	token, err = jwt.ParseString(auth, opts...)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getAuthSchema(c *fiber.Ctx) AuthSchema {
	auth := c.Get(fiber.HeaderAuthorization)
	if auth == "" {
		return AuthSchemaNone
	}
	if strings.HasPrefix(auth, "Basic") {
		return AuthSchemaBasic
	} else if strings.HasPrefix(auth, "Bearer") {
		return AuthSchemaBearer
	}
	return AuthSchemaOther
}
