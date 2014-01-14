package data

import (
	"fmt"
	"github.com/jbenet/s3"
	"github.com/jbenet/s3/s3util"
	"io"
	"strings"
)

type S3Store struct {
	bucket string
	domain string
	config *s3util.Config

	// used for auth credentials
	dataIndex *DataIndex
}

// format from `aws sts` cmd
type AwsCredentials struct {
	SecretAccessKey string
	SessionToken    string
	AccessKeyId     string
}

func NewS3Store(bucket string, index *DataIndex) (*S3Store, error) {

	if len(bucket) < 1 {
		return nil, fmt.Errorf("Invalid (empty) S3 Bucket name.")
	}

	if index == nil {
		return nil, fmt.Errorf("Invalid (nil) DataIndex.")
	}

	s := &S3Store{
		bucket:    bucket,
		domain:    "s3.amazonaws.com",
		dataIndex: index,
	}

	s.config = &s3util.Config{
		Service: s3.DefaultService,
		Keys:    new(s3.Keys),
	}

	return s, nil
}

func (s *S3Store) SetAwsCredentials(c *AwsCredentials) {
	s.config.AccessKey = c.AccessKeyId
	s.config.SecretKey = c.SecretAccessKey
	s.config.SecurityToken = c.SessionToken

	// pOut("Got Aws Credentials:\n")
	// pOut("	AccessKey: %s\n", s.config.AccessKey)
	// pOut("	SecretKey: %s\n", s.config.SecretKey)
	// pOut("	SessToken: %s\n\n", s.config.SecurityToken)
}

func (s *S3Store) AwsCredentials() *AwsCredentials {
	if s.config == nil || len(s.config.AccessKey) == 0 {
		return nil
	}

	return &AwsCredentials{
		AccessKeyId:     s.config.AccessKey,
		SecretAccessKey: s.config.SecretKey,
		SessionToken:    s.config.SecurityToken,
	}
}

func (s *S3Store) Url(key string) string {
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}
	return fmt.Sprintf("http://%s.%s%s", s.bucket, s.domain, key)
}

func (s *S3Store) Has(key string) (bool, error) {
	url := s.Url(key)
	rc, err := s3util.Open(url, s.config)

	if err == nil {
		rc.Close()
		return true, nil
	}

	if strings.Contains(err.Error(), "unwanted http status 404:") {
		return false, nil
	}

	return false, err
}

func (s *S3Store) Put(key string, value io.Reader) error {
	err := s.ensureUserAwsCredentials()
	if err != nil {
		return fmt.Errorf("aws credentials error: %v", err)
	}

	url := s.Url(key)
	w, err := s3util.Create(url, nil, s.config)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, value)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Store) Get(key string) (io.ReadCloser, error) {
	url := s.Url(key)
	return s3util.Open(url, s.config)
}

func (s *S3Store) getUserAwsCredentials() error {
	u := configUser()
	if !isNamedUser(u) {
		return fmt.Errorf("must be signed in to request aws credentials")
	}

	ui := s.dataIndex.NewUserIndex(u)
	c, err := ui.AwsCred()
	if err != nil {
		return err
	}

	s.SetAwsCredentials(c)
	return nil
}

func (s *S3Store) ensureUserAwsCredentials() error {
	// if we already have credentials, do nothing.
	if s.AwsCredentials() != nil {
		return nil
	}

	return s.getUserAwsCredentials()
}
