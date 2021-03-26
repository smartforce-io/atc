package githubservice

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
)

type AtcSettings struct {
	File   string `json:"fileDefault"`
	Prefix string `json:"prefix"`
	Field  string `json:"field"`
}

const (
	fileDefault   = "pom.xml"
	prefixDefault = ""
	fieldDefault  = "version"
)

var (
	errFailedResponse = errors.New("failed response")
	errNoAtcFile      = errors.New("no .atc.yaml file")
)

func getDefaultAtcSettings() *AtcSettings {
	return &AtcSettings{
		File:   fileDefault,
		Prefix: prefixDefault,
		Field:  fieldDefault,
	}
}

func getAtcSetting(client *github.Client, owner, repo string) (*AtcSettings, error) {
	in, resp, err := client.Repositories.DownloadContents(context.Background(), owner, repo, ".atc.yaml", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("getAtcSetting received a failed status for %s/%s: %q", owner, repo, resp.Status)
		return nil, errFailedResponse
	}
	defer in.Close()
	content, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	settings := &AtcSettings{}
	if err := yaml.Unmarshal(content, settings); err != nil {
		return nil, err
	}
	return settings, nil
}
