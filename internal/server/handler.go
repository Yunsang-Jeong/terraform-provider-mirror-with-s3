package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/internal/aws"
)

type WrapResponseWriter struct {
	w io.Writer
}

func (fw WrapResponseWriter) WriteAt(p []byte, offset int64) (n int, err error) {
	return fw.w.Write(p)
}

func (s *providerMirrorServer) setHandler(mux *http.ServeMux) error {
	for _, spec := range s.providerSpecs {
		index_url := fmt.Sprintf("/%s/index.json", spec.address)
		mux.HandleFunc(index_url, spec.listAvailableVersions)
		
		for _, v := range spec.versions{
			version_url := fmt.Sprintf("/%s/%s.json", spec.address, v.version)
			mux.HandleFunc(version_url, spec.listAvailableInstallationPackages)
		
			provider_url := fmt.Sprintf("/%s_%s_%s_%s.zip", v.name, v.version, v.os, v.arch)
			mux.HandleFunc(provider_url, spec.downloadPackage)
			fmt.Println(provider_url)
		}
	}

	return nil
}

// list available versions
// https://developer.hashicorp.com/terraform/internals/provider-network-mirror-protocol#list-available-versions
func (p *providerSpec) listAvailableVersions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := make(map[string]any)

	versions := make(map[string]any)
	for _, v := range p.versions {
		versions[v.version] = map[string]any{}
	}

	payload["versions"] = versions

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}


// list available installation packages
// https://developer.hashicorp.com/terraform/internals/provider-network-mirror-protocol#list-available-installation-packages
func (p *providerSpec) listAvailableInstallationPackages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := make(map[string]any)

	archives := make(map[string]any)
	last_path_segment := path.Base(r.URL.Path)
	version := strings.TrimSuffix(last_path_segment, filepath.Ext(last_path_segment))

	for _, v := range p.versions {
		if v.version == version {
			os_arch := strings.Join([]string{v.os, v.arch}, "_")
			archives[os_arch] =  map[string]any{
					"url":   fmt.Sprintf("%s_%s_%s_%s.zip", v.name, v.version, v.os, v.arch),
					"hahes": []string{}, // optional
			}
			break
		}
	}

	payload["archives"] = archives

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func (p *providerSpec) downloadPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")

	awsConfig, err := aws.NewAWSConfig("ap-northeast-2")
	if err != nil {
		return
	}

	awsConfig.DownloadObjectToBuffer(WrapResponseWriter{w: w}, p.bucketName, p.address + r.URL.Path)
}



