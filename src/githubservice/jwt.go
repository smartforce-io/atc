package githubservice

import (
	"crypto/x509"
	"encoding/pem"
	"envvars"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
	"time"
)

var (
	errNoPemEnv = errors.New("path to .pem is empty")
	errNoPem = errors.New("no .pem file")
)

func getJwt() (string, error) {
	pemPath := os.Getenv(envvars.PemPathVariable)
	if pemPath == "" { return "", errNoPemEnv}
	pemData, err := ioutil.ReadFile(pemPath)
	if err != nil { return "", err }
	block, _ := pem.Decode(pemData)
	if block == nil { return "", errNoPem }
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil { return "", err }
	atClaims := jwt.MapClaims{}
	atClaims["iat"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Unix()+(10 * 60)
	atClaims["iss"] = os.Getenv(envvars.AppId)
	at := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)
	token, err := at.SignedString(priv)
	if err != nil { return "", err}
	return token, nil
}