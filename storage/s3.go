package storage

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"xlsx-processor/pkg/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Downloading the file as bytes from S3
func DownloadFromS3Input(input types.Input) (fileContents []byte, err error) {
	fmt.Println("Downloading file from S3")
	fmt.Println(input)
	var reference types.SourceReference = input.Reference
	var credential types.Credential = input.Credential
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(reference.Region),
		Credentials: credentials.NewStaticCredentials(
			credential.Resources.Id,
			credential.Secrets.Secret,
			"",
		),
	})
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess)
	downloadInput := &s3.GetObjectInput{
		Bucket: aws.String(reference.Bucket),
		Key:    aws.String(reference.Prefix),
	}

	downloadResult, err := s3Client.GetObject(downloadInput)
	if err != nil {
		return nil, err
	}
	defer downloadResult.Body.Close()

	fileContents, err = io.ReadAll(downloadResult.Body)
	if err != nil {
		return nil, err
	}

	return fileContents, nil
}

// Uploading the file to an output S3 bucket
func UploadToS3Output(output types.Output, fileContents []byte) error {
	var reference types.SourceReference = output.Reference
	var credential types.Credential = output.Credential

	// Check if fileContents is empty
	if len(fileContents) == 0 {
		return fmt.Errorf("file contents are empty, nothing to upload")
	}

	// Log the size of the file contents for debugging
	log.Printf("Uploading file to S3: %d bytes", len(fileContents))

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(reference.Region),
		Credentials: credentials.NewStaticCredentials(
			credential.Resources.Id,
			credential.Secrets.Secret,
			"",
		),
	})
	if err != nil {
		return err
	}

	s3Client := s3.New(sess)

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(reference.Bucket),
		Key:         aws.String(reference.Prefix),
		Body:        bytes.NewReader(fileContents),
		ContentType: aws.String("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"),
	}

	_, err = s3Client.PutObject(uploadInput)
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}

// Downloading the file as bytes from S3
func DownloadFromS3(s3Key string, s3Bucket string, localPath string) error {
	accessKeyId := os.Getenv("S3_ACCESS_KEY_ID")
	secretAccessKeyId := os.Getenv("S3_SECRET_ACCESS_KEY")
	region := os.Getenv("S3_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyId,
			secretAccessKeyId,
			"",
		),
	})
	if err != nil {
		return err
	}

	s3Client := s3.New(sess)

	downloadInput := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key),
	}

	downloadResult, err := s3Client.GetObject(downloadInput)
	if err != nil {
		return err
	}
	defer downloadResult.Body.Close()

	fileContents, err := io.ReadAll(downloadResult.Body)
	if err != nil {
		return err
	}

	// Write the contents to the specified local path
	err = os.WriteFile(localPath, fileContents, 0644)
	if err != nil {
		return err
	}

	return nil
}
