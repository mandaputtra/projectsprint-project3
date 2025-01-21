package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mandaputtra/projectsprint-projects2/libs/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

const (
	AWS_S3_REGION = "ap-southeast-1"
	AWS_S3_BUCKET = "projectsprint-bucket-public-read"
)

var (
	AWS_ACCESS_KEY_ID     = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
	JWT_SECRET            = os.Getenv("JWT_SECRET")
	PORT                  = os.Getenv("PORT")
)

func isValidFile(file *multipart.FileHeader) (bool, string) {
	allowedExtensions := map[string]bool{
		".jpeg": true,
		".jpg":  true,
		".png":  true,
	}
	const maxFileSize = 100 * 1024 // 100 KiB

	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		return false, "Invalid file extension"
	}

	if file.Size > maxFileSize {
		return false, "File size exceeds the 100KiB limit"
	}

	return true, ""
}

func main() {
	// Initialize Gin router
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	router.POST("/v1/file", utils.Authorization, func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unable to get file from request",
			})
			return
		}

		if valid, msg := isValidFile(file); !valid {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": msg,
			})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to open file",
			})
			return
		}
		defer src.Close()

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(AWS_S3_REGION),
			Credentials: credentials.NewStaticCredentials(
				AWS_ACCESS_KEY_ID,
				AWS_SECRET_ACCESS_KEY,
				"", // a token will be created when the session it's used.
			),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create AWS session",
			})
			return
		}

		uploader := s3manager.NewUploader(sess)

		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(file.Filename),
			Body:   src,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to upload file, %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"uri": result.Location,
		})
	})

	// Start the server
	if err := router.Run(":" + PORT); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
