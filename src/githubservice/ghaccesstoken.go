package githubservice

import (
	"context"
	"errors"
	"net/http"
)

var (
	errWrongCreateAccessTokenStatus = errors.New("wrong access status during create access token for installation (not 201)")
)

func getAccessToken(id int64) (string, error) {
	jwt, err := getJwt()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	client := getGithubClient(jwt, ctx)
	inst, resp, err := client.Apps.CreateInstallationToken(ctx, id, nil)
	if err != nil { return "", err }
	if resp.StatusCode !=  http.StatusCreated {
		return "", errWrongCreateAccessTokenStatus
	}

	return inst.GetToken(), nil
}
