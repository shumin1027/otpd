package otp

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
	"time"
)

var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)

// Generate a OTP Secret
func GenerateSecret(secretLength ...int) string {
	size := 20
	if len(secretLength) > 0 && secretLength[0] > 0 {
		size = secretLength[0]
	}
	secret := make([]byte, size)
	_, err := rand.Reader.Read(secret)
	if err != nil {
		return ""
	}
	return b32NoPadding.EncodeToString(secret)
}

// Creates a TOTP token using the current time.
func GeneratePassCode(secret string) string {
	passcode, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return ""
	}
	return passcode
}

// Generate a new TOTP Key.
func GenerateKey(sub string, sec ...string) *otp.Key {
	var secret []byte
	var err error

	if len(sec) > 0 && len(sec[0]) > 0 {
		secret, err = b32NoPadding.DecodeString(sec[0])
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "C-MOM",
		AccountName: sub,
		Secret:      secret,
	})
	if err != nil {
		return nil
	}
	return key
}

// Convert OTP key into a PNG.
func GenerateQRCode(key *otp.Key) string {
	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return ""
	}
	png.Encode(&buf, img)
	sourcestring := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	return sourcestring
}

// Validate a TOTP using the current time.
func Validate(passcode, secret string) bool {
	return totp.Validate(passcode, secret)
}
