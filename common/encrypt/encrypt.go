package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/zeromicro/go-zero/core/codec"
	"golang.org/x/crypto/bcrypt"
	"io"
	"strings"
)

const (
	passwordEncryptSeed = "(erp)@#$"
	mobileAesKey        = "5A2E746B08D846502F37A6E2D85D583B"
	idCardHexKey        = "8a1b2c3d4e5f60718293a4b5c6d7e8f90123456789abcdef0123456789abcdef"
	accountHexKey       = "7a5e3d8c1f2b6a4d9e0c1f3a5b8d2e4f7a6c5d9e0f1a2b3c4d5e6f7a8b9c0d1e"
)

// EncAccount 使用 AES-GCM 加密银行卡号
func EncAccount(account string) (string, error) {
	key, err := hex.DecodeString(accountHexKey)
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", errors.New("AES-256 需要 32 字节密钥")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(account), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecAccount 使用 AES-GCM 解密银行卡号
func DecAccount(encryptedAccount string) (string, error) {
	key, err := hex.DecodeString(accountHexKey)
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", errors.New("AES-256 需要 32 字节密钥")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedAccount)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncIDCard 使用 AES-GCM 加密身份证号
func EncIDCard(idCard string) (string, error) {
	key, err := hex.DecodeString(idCardHexKey)
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", errors.New("AES-256 需要 32 字节密钥")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(idCard), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecIDCard 使用 AES-GCM 解密身份证号
func DecIDCard(encryptedIDCard string) (string, error) {
	key, err := hex.DecodeString(idCardHexKey)
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", errors.New("AES-256 需要 32 字节密钥")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedIDCard)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// EncMobile AES-ECB 模式，再base64加密
func EncMobile(mobile string) (string, error) {
	data, err := codec.EcbEncrypt([]byte(mobileAesKey), []byte(mobile))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func DecMobile(mobile string) (string, error) {
	originalData, err := base64.StdEncoding.DecodeString(mobile)
	if err != nil {
		return "", err
	}
	data, err := codec.EcbDecrypt([]byte(mobileAesKey), originalData)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Deprecated: EncPassword 已废弃，MD5 不安全，请使用 HashPassword（bcrypt）代替。
func EncPassword(password string) string {
	return Md5Sum([]byte(strings.TrimSpace(password + passwordEncryptSeed)))
}

// Md5Sum MD5 加密数据
func Md5Sum(data []byte) string {
	return hex.EncodeToString(byte16ToBytes(md5.Sum(data)))
}

func byte16ToBytes(in [16]byte) []byte {
	tmp := make([]byte, 16)
	for _, value := range in {
		tmp = append(tmp, value)
	}
	return tmp[16:]
}
