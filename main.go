package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func createS3Client() (*s3.Client, error) {
	godotenv.Load()

	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_ACCESS_KEY")

	// Load the AWS configuration with specific access key and secret key
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
	}
	return s3.NewFromConfig(cfg), nil
}

func uploadFileToS3(keyName string, fileContent []byte) error {
	s3Client, err := createS3Client()
	if err != nil {
		return err
	}

	godotenv.Load()
	bucketName := os.Getenv("BUCKET_NAME")

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &keyName,
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String("application/octet-stream"),
	})

	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %v", err)
	}

	return nil
}

func randomFileName() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "./static/index.html")
	})

	r.GET("/api/posts", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "hello"})
	})

	r.POST("/api/posts", func(ctx *gin.Context) {
		file, _ := ctx.FormFile("file")
		text := ctx.PostForm("text")

		in, _ := file.Open()
		defer in.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(in)
		fileContent := buf.Bytes()

		keyName := randomFileName()

		err := uploadFileToS3(keyName, fileContent)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"msg": file, "text": text})
	})

	r.Run(":8090")
}
