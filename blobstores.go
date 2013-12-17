package data

import (
	"fmt"
	"github.com/kr/s3"
	"github.com/kr/s3/s3util"
	"io"
	"os"
)

type kvStore interface {
	Put(key string, value io.Reader) error
	Get(key string) (io.ReadCloser, error)
}

type S3Store struct {
	bucket string
	config *s3util.Config
}

func NewS3Store(bucket string) (*S3Store, error) {

	if len(bucket) < 1 {
		return nil, fmt.Errorf("Invalid (empty) S3 Bucket name.")
	}

	s := &S3Store{bucket: bucket}

	err := s.setupConfig()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *S3Store) setupConfig() error {
	s.config.Service = s3.DefaultService
	s.config.Keys = new(s3.Keys)

	// move keys to config. for now use env key.
	s.config.AccessKey = os.Getenv("S3_ACCESS_KEY")
	s.config.SecretKey = os.Getenv("S3_SECRET_KEY")

	if len(s.config.AccessKey) < 1 {
		return fmt.Errorf("no S3_ACCESS_KEY env variable provided.")
	}

	return nil
}

func (s *S3Store) Put(key string, value io.Reader) error {
	w, err := s3util.Create(key, nil, s.config)
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
	return s3util.Open(key, s.config)
}
