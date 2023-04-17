package uploader

import (
	"fmt"
	"sync"

	"github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/internal/aws"
	errors "github.com/pkg/errors"
)

type uploader struct {
	bucket      string
	providers   []string
	os          string
	archtecture string
}

func NewUploader(bucket string, providers []string, os string, archtecture string) uploader {
	return uploader{
		bucket:      bucket,
		providers:   providers,
		os:          os,
		archtecture: archtecture,
	}
}

func (u *uploader) Start() error {
	var wg sync.WaitGroup

	wg.Add(len(u.providers))
	for _, p := range u.providers {
		p := p
		go func() {
			defer wg.Done()

			meta, err := getProviderMeta(p, u.os, u.archtecture)
			if err != nil {
				errors.Wrapf(err, "fail to upload provider: %s", p)
				return
			}

			reader, err := getProviderFileReader(meta.Description, meta.Version, u.os, u.archtecture)
			if err != nil {
				errors.Wrapf(err, "fail to upload provider: %s", p)
				return
			}

			awsConfig, err := aws.NewAWSConfig("ap-northeast-2")
			if err != nil {
				errors.Wrapf(err, "fail to upload provider: %s", p)
				return
			}

			key := fmt.Sprintf("registry.terraform.io/%s/%s/%s_%s_%s_%s.zip", meta.Namespace, meta.Name, meta.Description, meta.Version, u.os, u.archtecture)
			if err := awsConfig.UploadObject(u.bucket, key, reader); err != nil {
				errors.Wrapf(err, "fail to upload provider: %s", p)
				return
			}
		}()
	}
	wg.Wait()

	return nil
}
