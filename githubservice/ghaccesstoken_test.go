package githubservice

import (
	"fmt"
	"os"
	"testing"

	"github.com/smartforce-io/atc/envvars"
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
	//clear changes env for next tests
	os.Setenv(envvars.PemData, "")
}

func TestErrorGetPemFromPemPathVariable(t *testing.T) {
	mockClientProviderPtr := DefaultMockClientProvider()

	var tests = []struct {
		envvarsPem  string
		expectedErr string
	}{
		{"", "path to .pem is empty"},
		{"atcTest.pem", "open atcTest.pem: no such file or directory"},
		{"../../atcTestEmpty.pem", "no .pem file"},
	}
	for _, test := range tests {
		os.Setenv(envvars.PemPathVariable, test.envvarsPem)
		_, err := getAccessToken(10, mockClientProviderPtr)

		if fmt.Sprint(err) != test.expectedErr {
			t.Errorf("No get err with PemPathVariable = \"%s\"\nexpectedErr: \"%v\", got: \"%v\"", test.envvarsPem, test.expectedErr, err)
		}
	}
	//clear changes env for next tests
	os.Setenv(envvars.PemPathVariable, "")

}
