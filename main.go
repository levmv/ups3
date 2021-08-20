package main

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/url"
	"os"
)

var (
	input   *bufio.Reader
	outPath string
)

const help = `usage: upS3 [input file] <output path>
  - [input file]  If omitted then stdin will be used 
  - <output path> Must be in format of s3://bucket/dir/key`

func main() {

	if len(os.Args) == 1 {
		fmt.Println(help)
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		input = bufio.NewReader(os.Stdin)
		outPath = os.Args[1]
	} else {
		file, err := os.Open(os.Args[1])
		input = bufio.NewReader(file)
		if err != nil {
			die("Unable to open file %q", err)
		}
		defer file.Close()

		outPath = os.Args[2]
	}

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("ru-central1"),
		Endpoint: aws.String("https://storage.yandexcloud.net"),
	})

	u, err := url.Parse(outPath)
	if err != nil {
		die("Unable to parse file path %q", err)
	}

	bucketName := u.Host
	keyName := u.Path

	if len(keyName) > 0 && keyName[0] == '/' {
		keyName = keyName[1:]
	}

	if len(keyName) == 0 {
		die("Incorrect output path")
	}

	uploader := s3manager.NewUploader(sess)

	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &keyName,
		Body:   input,
	}

	// Perform an upload.
	_, err = uploader.Upload(upParams)

	if err != nil {
		die("Unable to upload %q to %q, %v", keyName, bucketName, err)
	}
}

func die(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
