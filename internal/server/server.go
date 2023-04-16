package server

import (
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type provider struct {
	hostname    string
	namespace   string
	_type       string
	name        string
	os          string
	archtecture string
	versions    []string
}

type providerMirrorServer struct {
	bucket string
	cert   string
	key    string
	router *gin.Engine
}

func NewProviderMirrorServer(bucket string, cert string, key string, release bool, silent bool) providerMirrorServer {
	if release {
		gin.SetMode(gin.ReleaseMode)
	}

	if silent {
		gin.DefaultWriter = io.Discard
	}

	return providerMirrorServer{
		bucket: bucket,
		cert:   cert,
		key:    key,
		router: gin.Default(),
	}
}

func (s *providerMirrorServer) Start() error {
	providers, err := listProvidersFromS3(s.bucket)
	if err != nil {
		return err
	}

	if s.router == nil {
		return errors.New("router is not initialized")
	}

	for _, p := range providers {
		baseUrl := fmt.Sprintf("%s/%s/%s", p.hostname, p.namespace, p._type)
		g := s.router.Group(baseUrl)

		g.GET("/index.json", p.listAvailableVersions)

		for _, v := range p.versions {
			g.GET(fmt.Sprintf("/%s.json", v), p.listAvailableInstallationPackages)

			providerFile := fmt.Sprintf("%s_%s_%s_%s.zip", p.name, v, p.os, p.archtecture)
			g.GET(providerFile, func(c *gin.Context) {
				providerFile := path.Base(c.FullPath())

				buffer, err := downloadProviderFromS3(s.bucket, fmt.Sprintf("%s/%s", baseUrl, providerFile))
				if err != nil {
					c.AbortWithError(500, err)
				}

				c.Data(200, "application/octet-stream", buffer)
			})
		}
	}

	return s.router.RunTLS(":443", s.cert, s.key)
}

// list available versions
// https://developer.hashicorp.com/terraform/internals/provider-network-mirror-protocol#list-available-versions
func (p *provider) listAvailableVersions(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	versions := make(map[string]any)
	for _, v := range p.versions {
		versions[v] = map[string]any{}
	}

	c.JSON(200, gin.H{
		"versions": versions,
	})
}

// list available installation packages
// https://developer.hashicorp.com/terraform/internals/provider-network-mirror-protocol#list-available-installation-packages
func (p *provider) listAvailableInstallationPackages(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")

	last_path_segment := path.Base(c.Request.URL.Path)
	version := strings.TrimSuffix(last_path_segment, filepath.Ext(last_path_segment))
	os_arch := fmt.Sprintf("%s_%s", p.os, p.archtecture)

	c.JSON(200, gin.H{
		"archives": gin.H{
			os_arch: gin.H{
				"url":   fmt.Sprintf("terraform-provider-%s_%s_%s.zip", p._type, version, os_arch),
				"hahes": []string{}, // optional
			},
		},
	})
}
