package githubservice

import "github.com/google/go-github/github"

type AtcSettings struct {
	File string `json:"fileDefault"`
	Prefix string `json:"prefix"`
	Field string `json:"field"`
}

const (
	fileDefault = "pom.xml"
	prefixDefault      = ""
	fieldDefault       = "version"
	)

func getDefaultAtcSettings() *AtcSettings {
	return &AtcSettings{
		File:   fileDefault,
		Prefix: prefixDefault,
		Field:  fieldDefault,
	}
}

func getAtcSetting(client *github.Client) (*AtcSettings, error) {




	return nil, nil
}
