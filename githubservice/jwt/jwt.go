package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"time"

	"github.com/smartforce-io/atc/envvars"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoPemEnv = errors.New("path to .pem is empty")
	errNoPem    = errors.New("no .pem file")
)

func GetJwt(pemData []byte) (string, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return "", errNoPem
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	atClaims := jwt.MapClaims{}
	atClaims["iat"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Unix() + (10 * 60)
	atClaims["iss"] = os.Getenv(envvars.AppId)
	at := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)
	token, err := at.SignedString(priv)
	if err != nil {
		return "", err
	}
	return token, nil
}
