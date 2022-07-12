package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shumin1027/otpd/pkg/http"
	"github.com/shumin1027/otpd/pkg/otp"
)

// @Summary Ping
// @Description ping
// @Produce application/json
// @Tags ping
// @Router /ping [GET]
// @Success	200 string string "pong"
func Ping(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).SendString("pong")
}

// 生成一个OTP密钥
func GetOTPKeyByNmae(c *fiber.Ctx) error {
	name := c.Query("name")
	//name := c.Locals("username").(string)

	account, err := otp.Get(name)
	if err == nil && account != nil {
		return http.Success(c, account)
	}

	secret := otp.GenerateSecret()
	key := otp.GenerateKey(name, secret)
	qr := otp.GenerateQRCode(key)

	account = &otp.Account{
		OTP:    key.URL(),
		Name:   name,
		QRCode: qr,
	}

	err = account.Save()
	if err != nil {
		return http.Error(c, err)
	}

	return http.Success(c, account)
}

// 获取当前验证码
func GetPassCodeByNmae(c *fiber.Ctx) error {
	name := c.Query("name")
	//name := c.Locals("username").(string)

	account, err := otp.Get(name)
	if err != nil || account == nil {
		return http.Fail(c, "no valid account found", http.StatusBadRequest)
	}

	key, _ := account.Key()
	passcode := otp.GeneratePassCode(key.Secret())

	return http.Success(c, passcode)
}

// 校验验证码是否有效
func Validate(c *fiber.Ctx) error {
	name := c.Query("name")
	passcode := c.Query("passcode")
	if len(name) == 0 {
		return http.Fail(c, "the name cannot be empty", http.StatusBadRequest)
	}
	if len(passcode) == 0 {
		return http.Fail(c, "the passcode cannot be empty", http.StatusBadRequest)
	}

	account, err := otp.Get(name)
	if err != nil || account == nil {
		return http.Fail(c, "no valid account found", http.StatusBadRequest)
	}

	key, _ := account.Key()

	ok := otp.Validate(passcode, key.Secret())
	return http.Success(c, ok)
}
