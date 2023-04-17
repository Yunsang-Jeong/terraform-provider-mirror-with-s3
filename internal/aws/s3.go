package aws

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	errors "github.com/pkg/errors"
)

func (a *awsConfig) ListBucketObjectKeys(bucketName string) ([]string, error) {
	keys := []string{}

	client := s3.NewFromConfig(a.config)

	pages := s3.NewListObjectsV2Paginator(client,
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		},
	)

	for pages.HasMorePages() {
		page, err := pages.NextPage(context.TODO())
		if err != nil {
			return nil, errors.Wrap(err, "fail to paginate of ListObjectV2")
		}

		for _, obj := range page.Contents {
			if path.Ext(*obj.Key) == ".zip" {
				keys = append(keys, *obj.Key)
			}
		}
	}

	return keys, nil
}

func (a *awsConfig) headObject(bucketName string, objectKey string) (*s3.HeadObjectOutput, error) {
	client := s3.NewFromConfig(a.config)

	return client.HeadObject(context.TODO(),
		&s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		},
	)
}

func (a *awsConfig) DownloadObjectToBuffer(bucketName string, objectKey string) ([]byte, error) {
	client := s3.NewFromConfig(a.config)

	headObject, err := a.headObject(bucketName, objectKey)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to headObject on %s/%s", bucketName, objectKey)
	}

	buf := make([]byte, int(headObject.ContentLength))
	writer := manager.NewWriteAtBuffer(buf)

	downloader := manager.NewDownloader(client)

	if _, err := downloader.Download(context.TODO(),
		writer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		},
	); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("fail to download object(%s) from aws s3 bucket(%s)", objectKey, bucketName))
	}

	return writer.Bytes(), nil
}

func (a *awsConfig) UploadObject(bucketName string, objectKey string, reader *io.ReadCloser) error {
	defer (*reader).Close()

	client := s3.NewFromConfig(a.config)

	uploader := manager.NewUploader(client)

	if _, err := uploader.Upload(context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   *reader,
		},
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("fail to updateload object(%s) from aws s3 bucket(%s)", objectKey, bucketName))
	}

	return nil
}
