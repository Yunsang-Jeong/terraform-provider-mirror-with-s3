package server

import (
	"fmt"
	"net/http"
)

type providerMirrorServer struct {
	bucketName  string
	awsS3Client *awsS3Client
	certFile    string
	keyFile     string
}

func NewProviderMirrorServer(bucketName, certFile, keyFile string) providerMirrorServer {
	return providerMirrorServer{
		bucketName:  bucketName,
		awsS3Client: nil,
		certFile:    certFile,
		keyFile:     keyFile,
	}
}

func (s *providerMirrorServer) Start() error {
	mux := http.NewServeMux()

	s.setHandler(mux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", 3000),
		Handler: mux,
	}

	return server.ListenAndServeTLS(s.certFile, s.keyFile)
}
