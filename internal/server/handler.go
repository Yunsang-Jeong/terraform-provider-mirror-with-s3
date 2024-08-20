package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type availableProvider struct {
	indexPath  string
	pHostname  string
	pNamespace string
	pType      string
}

func (s *providerMirrorServer) setHandler(mux *http.ServeMux) error {
	awsS3Client, err := newAWSS3Client("ap-northeast-2")
	if err != nil {
		return err
	}
	s.awsS3Client = awsS3Client

	availableProviders, err := s.getAvailableProviders()
	if err != nil {
		return err
	}

	for _, provider := range availableProviders {
		mux.HandleFunc(
			fmt.Sprintf("/%s/", provider.indexPath),
			s.proxyAWSS3Object,
		)
	}

	return nil
}

func (s *providerMirrorServer) getAvailableProviders() ([]*availableProvider, error) {
	objectKeys, err := s.awsS3Client.awsS3ListObjects(s.bucketName)
	if err != nil {
		return nil, err
	}

	availableProviders := make([]*availableProvider, 0)

	for _, objectKey := range objectKeys {
		if filepath.Base(objectKey) != "index.json" {
			continue
		}

		indexPath := filepath.Dir(objectKey)
		pathSegments := strings.Split(objectKey, "/")

		if len(pathSegments) != 4 {
			log.Printf("[skip] Weired Object path: %s", objectKey)
			continue
		}

		availableProviders = append(availableProviders, &availableProvider{
			indexPath:  indexPath,
			pHostname:  pathSegments[0],
			pNamespace: pathSegments[1],
			pType:      pathSegments[2],
		})

		log.Printf("sucess to regist available-provider: %s", indexPath)
	}

	return availableProviders, nil
}

func (s *providerMirrorServer) proxyAWSS3Object(w http.ResponseWriter, r *http.Request) {
	objectKey := strings.TrimPrefix(r.URL.Path, "/")

	if err := s.awsS3Client.awsS3ProxyObjectWithChunk(s.bucketName, objectKey, w); err != nil {
		log.Printf("Error %v\n", err)
		http.Error(w, "Error during download", http.StatusInternalServerError)
	}
}
