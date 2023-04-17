package uploader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	errors "github.com/pkg/errors"
)

type providerMeta struct {
	Id          string      `json:"id"`
	Namespace   string      `json:"namespace"`
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Description string      `json:"description"`
	Others      interface{} `json:"-"`
}

func getProviderMeta(provider string, os string, archtecture string) (*providerMeta, error) {
	url := fmt.Sprintf("https://registry.terraform.io/v1/providers/%s", provider)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("fail to get provider meta: status code is %d", resp.StatusCode))
	}

	meta := providerMeta{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&meta); err != nil {
		return nil, err
	}

	return &meta, nil
}

func getProviderFileReader(description string, version string, os string, archtecture string) (*io.ReadCloser, error) {
	baseUrl := fmt.Sprintf("https://releases.hashicorp.com/%s/%s/", description, version)
	fileName := fmt.Sprintf("%s_%s_%s_%s.zip", description, version, os, archtecture)
	resp, err := http.Get(baseUrl + fileName)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("fail to get provider file: status code is %d", resp.StatusCode))
	}

	return &resp.Body, nil
}
