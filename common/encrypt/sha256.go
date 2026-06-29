package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
)

func GenerateFileHash(file multipart.File) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	_, err := file.Seek(0, 0) // 重置文件指针
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
