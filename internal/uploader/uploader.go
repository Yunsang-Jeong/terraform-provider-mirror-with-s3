package uploader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Yunsang-Jeong/terraform-provider-mirror-with-s3/internal/aws"
	errors "github.com/pkg/errors"
)

type uploader struct {
	bucket      string
	providerHostname string
	providerNamespace string
	providerType string
	providerOS string
	providerArch string
	providerVersion string
}


type providerMeta struct {
	Id          string      `json:"id"`
	Namespace   string      `json:"namespace"`
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Description string      `json:"description"`
	Others      interface{} `json:"-"`
}

func NewUploader(bucket string, providerHostname string, providerNamespace string, providerType string, providerOS string, providerArch string, providerVersion string) uploader {
	return uploader{
		bucket:      bucket,
		providerHostname: providerHostname,
		providerNamespace: providerNamespace,
		providerType: providerType,
		providerOS: providerOS,
		providerArch: providerArch,
		providerVersion: providerVersion,
	}
}


func fetchJSON(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error fetching JSON from %s: status code is %d", url, resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return nil, err
	}

	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
    fmt.Printf("Error unmarshalling JSON: %v\n", err)
    return nil, err
	}

	return data, nil
}

func downloadAndUploadProvider(providerDownloadURL string, providerFileName string, bucketName string, objectKey string) (error) {
	resp, err := http.Get(providerDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("fail to get provider meta: status code is %d", resp.StatusCode))
	}

	awsConfig, err := aws.NewAWSConfig("ap-northeast-2")
	if err != nil {
		return err
	}
	
	if err := awsConfig.UploadObject(bucketName, objectKey, &resp.Body); err != nil {
		return err
	}

	return nil
}


func (u *uploader) Start() error {
	wellKnownURL := fmt.Sprintf("https://%s/.well-known/terraform.json", u.providerHostname)
	wellKnownJSON, err := fetchJSON(wellKnownURL)
	if err != nil {
		return err
	}

	providerInfoURL := fmt.Sprintf("https://%s%s%s/%s", u.providerHostname, wellKnownJSON["providers.v1"].(string), u.providerNamespace, u.providerType)
	
	if u.providerVersion == "" {
		providerInfoJSON, err := fetchJSON(providerInfoURL)
		if err != nil {
			return err
		}

		u.providerVersion = providerInfoJSON["version"].(string)
	}

	providerVersionInfoURL := fmt.Sprintf("%s/%s/download/%s/%s", providerInfoURL, u.providerVersion, u.providerArch, u.providerOS)
	providerVersionInfoJSON, err := fetchJSON(providerVersionInfoURL)
	if err != nil {
		return err
	}
	
	providerDownloadURL := providerVersionInfoJSON["download_url"].(string)
	providerFileName := providerVersionInfoJSON["filename"].(string)
	objectKey := fmt.Sprintf("%s/%s/%s/%s", u.providerHostname, u.providerNamespace, u.providerType, providerFileName)
	
	if downloadAndUploadProvider(providerDownloadURL, providerFileName, u.bucket, objectKey) != nil {
		return err
	}

	return nil
}
