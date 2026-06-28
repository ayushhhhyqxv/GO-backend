package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	bucket := "s3forgolang"
	key := "go/main.html"
	filepath := "main.html"
	err:= uploadtos3(bucket,key,filepath)
	if err!=nil{
		log.Fatalf("Error while uploading: %v",err)
	}
	fmt.Println("Done! Uploading ")
}

func uploadtos3(bucket, key, filepath string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	_,err =client.PutObject(context.TODO(),&s3.PutObjectInput{
		Bucket : &bucket,
		Key : &key,
		Body : file,
	})

	return err
}
