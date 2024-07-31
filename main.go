package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"week1/models"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func createS3Client() (*s3.Client, error) {
	godotenv.Load()

	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_ACCESS_KEY")
	region := os.Getenv("BUCKET_REGION")

	// Load the AWS configuration with specific access key and secret key
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
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

//func getObjectUrl(key string) (string, error) {
//	s3Client, err := createS3Client()
//	if err != nil {
//		return "", err
//	}
//
//	godotenv.Load()
//	bucketName := os.Getenv("BUCKET_NAME")
//
//	presignClient := s3.NewPresignClient(s3Client)
//	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
//		Bucket: aws.String(bucketName),
//		Key:    aws.String(key),
//	}, func(opts *s3.PresignOptions) {
//		opts.Expires = time.Duration(3600 * int64(time.Second))
//	})
//	if err != nil {
//		return "", err
//	}
//
//	return request.URL, nil
//}

func AddUrl(data []models.Info) ([]models.Item, error) {
	var items []models.Item
	for _, info := range data {
		objectUrl := "http://d2hgyikbqxrz2l.cloudfront.net/" + info.Filename
		//objectUrl, err := getObjectUrl(info.Filename)
		//if err != nil {
		//	return nil, err
		//}
		items = append(items, models.Item{
			Info:      info,
			ObjectUrl: objectUrl,
		})
	}
	return items, nil
}

type Presigner struct {
	PresignClient *s3.PresignClient
}

func (presigner Presigner) GetObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

func randomFileName(origFileName string) string {
	ext := filepath.Ext(origFileName)
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b) + ext
}

func fileByteContent(file *multipart.FileHeader) []byte {
	in, _ := file.Open()
	defer in.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)
	fileContent := buf.Bytes()
	return fileContent
}

func main() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "./static/index.html")
	})

	r.GET("/loaderio-5402094e37c9a38d5a9c90af4c9fd4f5.txt", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "./static/loaderio-5402094e37c9a38d5a9c90af4c9fd4f5.txt")
	})

	r.GET("/api/posts", func(ctx *gin.Context) {
		rawData, err := models.GetData()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data, err := AddUrl(rawData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": data})
	})

	r.POST("/api/posts", func(ctx *gin.Context) {
		file, _ := ctx.FormFile("file")
		text := ctx.PostForm("text")

		fileName := file.Filename

		fileContent := fileByteContent(file)
		keyName := randomFileName(fileName)

		err := uploadFileToS3(keyName, fileContent)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = models.InsertData(text, keyName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"msg": file, "text": text})
	})

	r.Run(":8080")
}
