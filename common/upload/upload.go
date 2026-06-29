package upload

import (
	"context"
	"erp/common/encrypt"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

var (
	ContentTypeErr = errors.New("上传的图片类型非 jpg、jpeg、png")
)

type Oss interface {
	UploadFile(ctx context.Context, file multipart.File, filename string) (string, error)
}

func getContentTypeAndSha256Filename(file multipart.File, filename string) (contentType, sha256Filename string, err error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	default:
		contentType = "application/octet-stream"
	}

	// 改文件名为 sha256
	key, err := encrypt.GenerateFileHash(file)
	if err != nil {
		return "", "", err
	}
	key = key + ext
	return contentType, key, nil
}
