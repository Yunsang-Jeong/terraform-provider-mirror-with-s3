package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type awsS3Client struct {
	client    *s3.Client
	chunkSize int64
}

func newAWSS3Client(region string) (*awsS3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create a new aws client")
	}

	return &awsS3Client{
		client:    s3.NewFromConfig(cfg),
		chunkSize: int64(5 * 1024 * 1024),
	}, nil
}

func (c *awsS3Client) awsS3ListObjects(bucketName string) ([]string, error) {
	objectKeys := []string{}

	pages := s3.NewListObjectsV2Paginator(
		c.client,
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		},
	)

	for pages.HasMorePages() {
		page, err := pages.NextPage(context.TODO())
		if err != nil {
			return nil, errors.Wrap(err, "fail to paginate of ListObjectV2")
		}

		for _, object := range page.Contents {
			objectKeys = append(objectKeys, *object.Key)
		}
	}

	return objectKeys, nil
}

type awsS3ObjectInfo struct {
	contentLength int64
	contentType   string
}

func (c *awsS3Client) awsS3GetObjectInfo(bucketName string, objectKey string) (*awsS3ObjectInfo, error) {
	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	objectInfo, err := c.client.HeadObject(context.TODO(), headObjectInput)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to head-object: %s/%s", bucketName, objectKey)
	}

	return &awsS3ObjectInfo{
		contentLength: *objectInfo.ContentLength,
		contentType:   *objectInfo.ContentType,
	}, nil
}

type awsS3ObjectChunk struct {
	rangeStart int64
	data       *[]byte
	err        error
}

func (c *awsS3Client) awsS3ProxyObjectWithChunk(bucketName string, objectKey string, w http.ResponseWriter) error {
	info, err := c.awsS3GetObjectInfo(bucketName, objectKey)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", info.contentType)

	chunkCount := (info.contentLength + c.chunkSize - 1) / c.chunkSize
	chunks := make(chan awsS3ObjectChunk, chunkCount)

	var wg sync.WaitGroup

	for i := int64(0); i < chunkCount; i++ {
		rangeStart := i * c.chunkSize
		rangeEnd := rangeStart + c.chunkSize - 1
		if rangeEnd >= info.contentLength {
			rangeEnd = info.contentLength - 1
		}

		wg.Add(1)

		go func(rangeStart, rangeEnd int64) {
			defer wg.Done()

			object, err := c.client.GetObject(context.TODO(),
				&s3.GetObjectInput{
					Bucket: aws.String(bucketName),
					Key:    aws.String(objectKey),
					Range:  aws.String(fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd)),
				},
			)
			if err != nil {
				chunks <- awsS3ObjectChunk{
					err: err,
				}
				return
			}
			defer object.Body.Close()

			bufferSize := rangeEnd - rangeStart + 1
			buffer := make([]byte, bufferSize)
			_, err = io.ReadFull(object.Body, buffer)
			if err != nil && err != io.ErrUnexpectedEOF {
				chunks <- awsS3ObjectChunk{
					err: err,
				}
				return
			}

			chunks <- awsS3ObjectChunk{
				rangeStart: rangeStart,
				data:       &buffer,
			}
		}(rangeStart, rangeEnd)
	}

	go func() {
		wg.Wait()
		close(chunks)
	}()

	expectedRangeStart := int64(0)
	receivedChunks := make(map[int64]*[]byte)

	for chunk := range chunks {
		if chunk.err != nil {
			return chunk.err
		}

		if chunk.rangeStart == expectedRangeStart {
			_, err := w.Write(*chunk.data)
			if err != nil {
				return err
			}

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

			expectedRangeStart += int64(len(*chunk.data))

			for {
				if data, exists := receivedChunks[expectedRangeStart]; exists {
					_, err := w.Write(*data)
					if err != nil {
						log.Printf("Failed to send part to client: %v", err)
						return err
					}

					if f, ok := w.(http.Flusher); ok {
						f.Flush()
					}

					delete(receivedChunks, expectedRangeStart)
					expectedRangeStart += int64(len(*data))
				} else {
					break
				}
			}
		} else {
			receivedChunks[chunk.rangeStart] = chunk.data
		}
	}

	log.Printf("Successfully sent object: %s", objectKey)
	return nil
}
