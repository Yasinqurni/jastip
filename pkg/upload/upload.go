package upload

import (
	"context"
	"fmt"
	"jastip-app/config"
	"jastip-app/pkg/logger"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Media struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Path        string `json:"path"`
	Media       multipart.File
}

type UploadOutput struct {
	Path string
	Url  string
}

type Uploader struct {
	config    *config.Config
	s3Session *session.Session
	mtx       sync.Mutex
}

func (u *Uploader) ConnectS3(config *config.Config) {
	u.config = config
	s3session, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AWS.AccessKeyID),
		Credentials: credentials.NewStaticCredentials(config.AWS.AccessKeyID, config.AWS.SecretAcessKey, ""),
	})

	if err != nil {
		logger.L().Error(err.Error())
	}

	u.s3Session = s3session
}

func (u *Uploader) Upload(ctx context.Context, media *Media) (*UploadOutput, error) {
	u.mtx.Lock()
	defer u.mtx.Unlock()
	uploader := s3manager.NewUploader(u.s3Session)
	mime := strings.Replace(media.ContentType, "image/", "", 1)
	fullFileName := fmt.Sprintf("%s.%s", media.Filename, mime)
	fullPath := filepath.Join(media.Path, fullFileName)

	input := &s3manager.UploadInput{
		Bucket:      aws.String(u.config.AWS.Bucket), // bucket's name
		Key:         aws.String(fullPath),            // files destination location
		Body:        media.Media,                     // content of the media
		ContentType: aws.String(media.ContentType),   // content type
	}

	output, err := uploader.Upload(input)
	if err != nil {
		logger.L().Error(err.Error())
		return nil, err
	}

	return &UploadOutput{
		Url:  output.Location,
		Path: fullPath,
	}, nil
}

func (u *Uploader) Delete(ctx context.Context, path string) error {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	s3Client := s3.New(u.s3Session)
	iter := s3manager.NewDeleteListIterator(s3Client, &s3.ListObjectsInput{
		Bucket: aws.String(u.config.AWS.Bucket),
		Prefix: aws.String(path),
	})

	if err := s3manager.NewBatchDeleteWithClient(s3Client).Delete(ctx, iter); err != nil {
		logger.L().Error(err.Error())
		return err
	}

	return nil
}
