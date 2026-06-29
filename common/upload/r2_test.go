package upload

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestR2Upload(t *testing.T) {
	r2conf := &R2Conf{
		AccessKeyId:     "6b16230d08d9dfbd48c40033fa6ebdb8",
		AccessKeySecret: "f16efa9458918034479b703ae250697f23ab8776f2bdc7660322d622b7df2b1b",
		Endpoint:        "https://7d97f876a089983c3602465d501ab460.r2.cloudflarestorage.com",
		BucketName:      "cjsakana",
		Domain:          "sakanatang.dpdns.org",
	}
	OssClient := NewR2Client(r2conf)

	filePath := "D:\\瞅瞅\\ygo\\5F3571E1A2FCEBE7CC592B6380B79808.jpg"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("无法打开文件 %q: %v", filePath, err)
	}
	defer file.Close()

	url, err := OssClient.UploadFile(context.Background(), file, "5F3571E1A2FCEBE7CC592B6380B79808.jpg")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(url)
}
