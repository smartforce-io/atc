package accesstoken

import (
	"fmt"
	"os"
	"testing"

	"github.com/smartforce-io/atc/envvars"
	"github.com/smartforce-io/atc/githubservice/provider"
)

const (
	expected = "aaa"

	testRsaKey = `
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`
)

func TestGetAccessTokenBasic(t *testing.T) {
	mockClientProviderPtr := provider.DefaultMockClientProvider()
	envvarsPemDataOld := envvars.PemData
	os.Setenv(envvars.PemData, testRsaKey)
	defer os.Setenv(envvars.PemData, envvarsPemDataOld)

	token, err := GetAccessToken(10, mockClientProviderPtr)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	if token != expected {
		t.Errorf("Unexpected token, expected %s, got %s", expected, token)
	}
}

func TestErrorGetPemFromPemPathVariable(t *testing.T) {
	mockClientProviderPtr := provider.DefaultMockClientProvider()
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
		_, err := GetAccessToken(10, mockClientProviderPtr)

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
