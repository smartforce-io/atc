package accesstoken

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/smartforce-io/atc/githubservice/provider"

	"github.com/smartforce-io/atc/envvars"
	"github.com/smartforce-io/atc/githubservice/jwt"
)

var (
	errWrongCreateAccessTokenStatus = errors.New("wrong access status during create access token for installation (not 201)")
)

func GetAccessToken(id int64, clientProvider provider.ClientProvider) (string, error) {
	var pemData []byte
	var err error
	pemEnv := os.Getenv(envvars.PemData)
	if pemEnv == "" {
		pemPath := os.Getenv(envvars.PemPathVariable)
		if pemPath == "" {
			return "", jwt.ErrNoPemEnv
		}
		pemData, err = os.ReadFile(pemPath)
		if err != nil {
			return "", err
		}
		log.Printf("ATC uses pem from file: %q", pemPath)

	} else {
		pemData = []byte(pemEnv)
		log.Print("ATC uses pem data from environment variable")
	}

	j, err := jwt.GetJwt(pemData)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	client := clientProvider.Get(j, ctx)
	inst, resp, err := client.Apps.CreateInstallationToken(ctx, id, nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusCreated {
		return "", errWrongCreateAccessTokenStatus
	}

	return inst.GetToken(), nil
}
