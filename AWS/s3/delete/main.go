package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	bucket := "s3forgolang"
	key := "go/golang.txt"
	cfg, err := config.LoadDefaultConfig(context.TODO())

	client := s3.NewFromConfig(cfg)

	_, err = client.DeleteObject(context.TODO(),&s3.DeleteObjectInput{
		Bucket: &bucket,
		Key: &key,
	})

	if err!=nil {
		panic(err)
	}

	fmt.Println("Deleted Object -> ",key)
}
