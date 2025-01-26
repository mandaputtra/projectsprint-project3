package handlers

import (
	"fmt"
	"net/http"
	"project3/services/ms-upp-svc/database"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

// Service
func (a *APIEnv) UploadFile(c *gin.Context) {
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
		Region: aws.String(a.ENV.AWS_S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			a.ENV.AWS_ACCESS_KEY_ID,
			a.ENV.AWS_SECRET_ACCESS_KEY,
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
		Bucket: aws.String(a.ENV.AWS_S3_BUCKET),
		Key:    aws.String(file.Filename),
		Body:   src,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to upload file, %v", err),
		})
		return
	}

	// Insert to database
	db := a.DB
	insertedFile := database.File{
		FileUri:          result.Location,
		FileThumbnailUri: result.Location,
	}
	db.Create(&insertedFile)

	c.JSON(http.StatusOK, gin.H{
		"fileId":           insertedFile.FileID,
		"fileUri":          result.Location,
		"FileThumbnailUri": result.Location,
	})
}
