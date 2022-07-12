package http

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	self "github.com/shumin1027/otpd/app"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	_ "github.com/shumin1027/otpd/docs"
	"github.com/shumin1027/otpd/http/middleware/auth"
	"github.com/shumin1027/otpd/pkg/http"
	log "github.com/shumin1027/otpd/pkg/logger"
	"go.uber.org/zap"
)

var app = fiber.New(fiber.Config{
	AppName: self.Name,
})

const SECRET = "a72cd591-5b57-4ffe-b6e3-e99c317ff43c"

var SecretKey jwk.Key

func init() {
	key, err := jwk.New([]byte(SECRET))
	if err != nil {
		panic(any("error for gen secret key"))
	}
	SecretKey = key
}

// OTP Server API
// @title OTP Server API
// @version 1.0
// @Description OTP Server API
// @host localhost:18181
// @BasePath /
func Start(addr string) {
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:       "${time} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat:   "2006/01/02 15:04:05",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
	}))

	//app.Use(authentication())

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/stack", func(c *fiber.Ctx) error {
		return c.JSON(app.Stack())
	})

	app.Get("/ping", Ping)
	app.Get("/key", GetOTPKeyByNmae)
	app.Get("/validate", Validate)
	app.Get("/passcode", GetPassCodeByNmae)

	go func() {
		// service connections
		log.S().Fatal(app.Listen(addr))
	}()

	graceful(app)
}

func authentication() fiber.Handler {
	return auth.New(auth.Config{
		ContextKey:      "auth",
		SigningKey:      SecretKey,
		JWTParseOptions: []jwt.ParseOption{jwt.WithSubject("AccessToken")},
		AuthSchemas:     []auth.AuthSchema{auth.AuthSchemaBasic, auth.AuthSchemaBearer},
		Filter: func(c *fiber.Ctx) bool {
			//在验证token之前会调用此方法,如果返回true,则不验证token,直接调用接口
			uri := string(c.Request().RequestURI())
			if strings.HasPrefix(uri, "/validate") || strings.HasPrefix(uri, "/swagger") || uri == "/ping" || uri == "/version" || uri == "/stack" {
				return true
			} else {
				return false
			}
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			//如果token验证失败,会调用此方法,err包含Token验证失败的原因,例如过期
			log.L().Error("token validation failed", zap.Error(err))
			return c.Status(http.StatusUnauthorized).SendString(err.Error())
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			result := c.Locals("auth").(*auth.AuthResult)
			token := result.Token
			switch result.AuthSchema {
			case auth.AuthSchemaBasic:
				{
					c.Locals("username", result.User.Username)
				}
			case auth.AuthSchemaBearer:
				{
					if username, ok := token.Get("username"); ok {
						c.Locals("username", username.(string))
					}
				}
			}
			return c.Next()
		},
	})
}

func graceful(app *fiber.App) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.S().Info("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		log.S().Fatal(fmt.Sprintf("server shutdown:%s", err.Error()))
	}

	log.S().Info("running cleanup tasks ...")
	// Your cleanup tasks go here

	<-ctx.Done()
	log.S().Info("server exiting")
}
