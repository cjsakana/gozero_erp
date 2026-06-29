package upload

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"mime/multipart"
)

var (
	CreateR2ClientFail = errors.New("连接R2失败")
	UploadFileFail     = errors.New("上传文件失败")
)

type R2Conf struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	BucketName      string
	Domain          string
}

type R2Client struct {
	*s3.Client
	BucketName string
	Domain     string
}

func NewR2Client(conf *R2Conf) Oss {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AccessKeyId, conf.AccessKeySecret, "")),
		config.WithRegion("auto"), // Cloudflare R2 推荐 "auto"
		//config.WithHTTPClient(&http.Client{
		//	Timeout: 30 * time.Second, // 单次请求超时
		//}),
		//config.WithRetryer(func() aws.Retryer {
		//	return retry.NewStandard(func(o *retry.StandardOptions) {
		//		o.MaxAttempts = 5 // 增加重试次数（默认3次）
		//	})
		//}),
	)
	if err != nil {
		log.Fatal(CreateR2ClientFail)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(conf.Endpoint)
	})
	return &R2Client{
		Client:     client,
		BucketName: conf.BucketName,
		Domain:     conf.Domain,
	}
}

func (m *R2Client) UploadFile(ctx context.Context, file multipart.File, filename string) (string, error) {
	contentType, sha256Filename, err := getContentTypeAndSha256Filename(file, filename)
	if err != nil {
		return "", err
	}

	// 上传对象
	_, err = m.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(m.BucketName),
		Key:         aws.String(sha256Filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", UploadFileFail
	}
	url := "https://" + m.Domain + "/" + sha256Filename
	return url, nil
}
