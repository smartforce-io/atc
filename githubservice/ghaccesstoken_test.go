package githubservice

import (
	"fmt"
	"log"
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
	envvarsPemDataOld := envvars.PemData
	os.Setenv(envvars.PemData, testRsaKey)
	defer os.Setenv(envvars.PemData, envvarsPemDataOld)

	token, err := getAccessToken(10, mockClientProviderPtr)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	if token != expected {
		t.Errorf("Unexpected token, expected %s, got %s", expected, token)
	}
}

func TestErrorGetPemFromPemPathVariable(t *testing.T) {
	mockClientProviderPtr := DefaultMockClientProvider()
	envvarsPemDataOld := envvars.PemData
	envvarsPemPathOld := envvars.PemPathVariable
	os.Setenv(envvars.PemData, "")
	defer os.Setenv(envvars.PemData, envvarsPemDataOld)
	defer os.Setenv(envvars.PemPathVariable, envvarsPemPathOld)

	//create atcTestEmpty.pem for test
	file, err := os.Create("atcTestEmpty.pem")
	if err != nil {
		t.Error(err)
	}
	file.Close()

	var tests = []struct {
		envvarsPem  string
		expectedErr string
	}{
		{"", "path to .pem is empty"},
		{"atcTest.pem", "open atcTest.pem: no such file or directory"},
		{"atcTestEmpty.pem", "no .pem file"},
	}
	for _, test := range tests {
		os.Setenv(envvars.PemPathVariable, test.envvarsPem)
		_, err := getAccessToken(10, mockClientProviderPtr)

		if fmt.Sprint(err) != test.expectedErr {
			t.Errorf("No get err with PemPathVariable = \"%s\"\nexpectedErr: \"%v\", got: \"%v\"", test.envvarsPem, test.expectedErr, err)
		}
	}
	//del atcTestEmpty.pem for this test
	err = os.Remove("atcTestEmpty.pem")
	if err != nil {
		t.Error(err)
	}
}
