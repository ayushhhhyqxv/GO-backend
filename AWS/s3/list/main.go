package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)


func main(){
	bucket:="s3forgolang"
	listfiles(bucket)
}

func listfiles(bucket string){
	cfg,_:= config.LoadDefaultConfig(context.TODO())

	client:= s3.NewFromConfig(cfg)

	resp,_ := client.ListObjectsV2(context.TODO(),&s3.ListObjectsV2Input{
		Bucket: &bucket,
	})

	for _,item:=range resp.Contents {
		fmt.Printf("\nContent -> %v , Size of content-> %v bytes\n",*item.Key,*item.Size)
	}
}