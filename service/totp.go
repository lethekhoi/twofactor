package service

import (
	"bytes"
	"image/png"

	"github.com/lethekhoi/twofactor/entities"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func TOTPKey(user *entities.User) (*otp.Key, bytes.Buffer, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Secret:      []byte("MySecretKey"),
		Issuer:      "Login",
		AccountName: user.Username,
	})
	if err != nil {
		panic(err)
	}
	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)

	return key, buf, nil

}
