package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type PluginSettings struct {
	AccessKeyId string                `json:"accessKeyId"`
	Region      string                `json:"endpoint"`
	Secrets     *SecretPluginSettings `json:"-"`
	InstanceId  string                `json:"instanceId"`
}

type SecretPluginSettings struct {
	AccessKeySecret string `json:"accessKeySecret"`
}

func LoadPluginSettings(source backend.DataSourceInstanceSettings) (*PluginSettings, error) {
	if source.JSONData == nil || len(source.JSONData) < 1 {
		return nil, errors.New("AccessKeySecret cannot be null")
	}

	settings := PluginSettings{}

	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal PluginSettings json: %w", err)
	}

	settings.Secrets = loadSecretPluginSettings(source.DecryptedSecureJSONData)

	return &settings, nil
}

func loadSecretPluginSettings(source map[string]string) *SecretPluginSettings {
	return &SecretPluginSettings{
		AccessKeySecret: source["accessKeySecret"],
	}
}
