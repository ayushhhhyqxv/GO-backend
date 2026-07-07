package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	bucket := "s3forgolang"
	key := "go/main.html"
	filepath := "websockets.html"
	err := downloadfroms3(bucket, key, filepath)
	if err != nil {
		log.Fatalf("Error while Downloading %v", err)
	}
	fmt.Println("Downloaded required files")
}

func downloadfroms3(bucket, key, filepath string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	return err
}
