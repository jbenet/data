package data

import (
	"fmt"
	"github.com/kr/s3"
	"github.com/kr/s3/s3util"
	"os"
	"strings"
)

const datadexS3Bucket = "datadex.archives"

func uploadCmd(args []string) error {
	if len(args) < 1 {
		return UploadDataset("datadex")
	} else {
		return UploadDataset(args[0])
	}
}

func UploadDataset(service string) error {

	// ensure the dataset has required information
	err := fillOutDatafileInPath(DatasetFile)
	if err != nil {
		return err
	}

	switch strings.ToLower(service) {
	case "datadex":
		return uploadDatasetToDatadex()
	}

	return fmt.Errorf("Unsupported storage service %s", service)
}

func uploadDatasetToDatadex() error {

	s3config, err := s3config()
	if err != nil {
		return err
	}

	return uploadDatasetToS3(s3config, datadexS3Bucket)
}

func uploadDatasetToS3(config *s3util.Config, bucket string) error {
	return fmt.Errorf("uploadDatasetToS3 not implemented")
}

func s3config() (*s3util.Config, error) {
	c := &s3util.Config{
		Service: s3.DefaultService,
		Keys:    new(s3.Keys),
	}

	// move keys to config. for now use env key.
	c.AccessKey = os.Getenv("S3_ACCESS_KEY")
	c.SecretKey = os.Getenv("S3_SECRET_KEY")

	if len(c.AccessKey) < 1 {
		return nil, fmt.Errorf("no S3_ACCESS_KEY env variable provided.")
	}

	return c, nil
}
