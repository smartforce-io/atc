package githubservice

import (
	"envvars"
	"os"
	"testing"
)

const (
	tokenResponse = `
{
	"token" : "aaa",
	"expires_at": "2016-07-11T22:14:10Z"
}
`
	expected = "aaa"
)

func TestGetAccessTokenBasic(t *testing.T) {
	mockClientProviderPtr := DefaultMockClientProvider()

	os.Setenv(envvars.PemData, testRsaKey)

	token, err := getAccessToken(10, mockClientProviderPtr)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	if token != expected {
		t.Errorf("Unexpected token, expected %s, got %s", expected, token)

	}
}
