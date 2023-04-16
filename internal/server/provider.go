package server

import (
	"path/filepath"
	"strings"

	"github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/internal/aws"
)

func listProvidersFromS3(bucket string) ([]provider, error) {
	awsConfig, err := aws.NewAWSConfig("ap-northeast-2")
	if err != nil {
		return nil, err
	}

	objectKeys, err := awsConfig.ListBucketObjectKeys(bucket)
	if err != nil {
		return nil, err
	}

	providers := map[string]provider{}
	for _, key := range objectKeys {
		pathSegments := strings.Split(key, "/")
		if len(pathSegments) != 4 {
			continue
		}

		hostname := pathSegments[0]
		namespace := pathSegments[1]
		_type := pathSegments[2]
		fileName := pathSegments[3]

		fileNameSegments := strings.Split(strings.TrimSuffix(fileName, filepath.Ext(fileName)), "_")
		if len(fileNameSegments) != 4 {
			continue
		}

		providerName := fileNameSegments[0]
		providerVersion := fileNameSegments[1]
		providerOS := fileNameSegments[2]
		providerArchtecture := fileNameSegments[3]

		key := strings.Join([]string{hostname, namespace, _type}, "/")

		if p, exists := providers[key]; exists {
			p.versions = append(p.versions, providerVersion)
			providers[key] = p
		} else {
			providers[key] = provider{
				hostname:    hostname,
				namespace:   namespace,
				_type:       _type,
				name:        providerName,
				os:          providerOS,
				archtecture: providerArchtecture,
				versions:    []string{providerVersion},
			}
		}
	}

	values := make([]provider, 0, len(providers))
	for _, v := range providers {
		values = append(values, v)
	}

	return values, nil
}

func downloadProviderFromS3(bucket string, object string) ([]byte, error) {
	awsConfig, err := aws.NewAWSConfig("ap-northeast-2")
	if err != nil {
		return nil, err
	}

	buffer, err := awsConfig.DownloadObjectToBuffer(bucket, object)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
