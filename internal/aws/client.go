package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	errors "github.com/pkg/errors"
)

type awsConfig struct {
	config aws.Config
}

func NewAWSConfig(region string) (*awsConfig, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create a new aws client")
	}

	return &awsConfig{
		config: cfg,
	}, nil
}
