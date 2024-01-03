package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/internal/aws"
)


type providerMirrorServer struct {
	bucketName string
	providerSpecs []*providerSpec
}

type providerSpec struct {
	bucketName string
	address string
	versions []*providerVersionSpec
}

type providerVersionSpec struct {
	name string
	version string
	os string
	arch string
}

func NewProviderMirrorServer(bucketName string) providerMirrorServer {
	return providerMirrorServer{
		bucketName: bucketName,
		providerSpecs: nil,
	}
}

func (s *providerMirrorServer) Start() error {
	if err := s.getAvailableProvidersFromS3Bucket() ;err != nil {
		return err
	}

	mux := http.NewServeMux()
	
	s.setHandler(mux)

	server := http.Server{
		Addr: fmt.Sprintf(":%d", 3000),
		Handler: mux,
	}

	generateCerticiate("key.pem", "cert.pem")

	return server.ListenAndServeTLS("cert.pem", "key.pem")
}

func (s *providerMirrorServer) getAvailableProvidersFromS3Bucket() error {
	awsConfig, err := aws.NewAWSConfig("ap-northeast-2")
	if err != nil {
		return err
	}

	objectKeys, err := awsConfig.ListBucketObjectKeys(s.bucketName)
	if err != nil {
		return err
	}

	providers := map[string][]*providerVersionSpec{}

	for _, key := range objectKeys {
		pathSegments := strings.Split(key, "/")
		if len(pathSegments) != 4 {
			continue
		}

		hostnamePortion := pathSegments[0]
		namespacePortion := pathSegments[1]
		typePortion := pathSegments[2]
		providerAddress := strings.Join([]string{hostnamePortion, namespacePortion, typePortion}, "/")

		providerFileName := pathSegments[3]
		fileNameSegments := strings.SplitN(strings.TrimSuffix(providerFileName, filepath.Ext(providerFileName)), "_", 4)
		if len(fileNameSegments) != 4 {
			continue
		}
		providerName := fileNameSegments[0]
		providerVersion := fileNameSegments[1]
		providerOS := fileNameSegments[2]
		providerArch := fileNameSegments[3]

		spec := &providerVersionSpec{
			name: providerName,
			version: providerVersion,
			os: providerOS,
			arch: providerArch,
		}

		if p, exists := providers[providerAddress]; exists {
			providers[providerAddress] = append(p, spec)
		} else {
			providers[providerAddress] = []*providerVersionSpec{spec}
		}
	}

	s.providerSpecs = make([]*providerSpec, 0, len(providers))
	for address, specs := range providers {
		s.providerSpecs = append(s.providerSpecs, &providerSpec{
			bucketName: s.bucketName,
			address: address,
			versions: specs,
		})
	}

	return nil
}
