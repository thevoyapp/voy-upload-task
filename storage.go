package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/aws/session"
  "os"
  "bytes"
  "fmt"
  "github.com/google/uuid"
  "strings"
  "errors"
  "net/url"
)

const (
  OS_REGION = "S3_REGION"
  OS_BUCKET = "S3_BUCKET"
  DEFAULT_REGION = "us-east-1"
  DEFAULT_BUCKET = "static.thevoyapp.com"
)

func getValue(defaultVal string, osString string) string {
  value, ok := os.LookupEnv(osString)
  if ok {return value} else {return defaultVal}
}

func FileKey(path string, contentType string) (string, error) {
  splitContent := strings.Split(contentType, "/")
  if len(splitContent) != 2 {
    return "", errors.New("Invalid Content Type " + contentType)
  }
  return fmt.Sprintf("%s/%s.%s", path, uuid.New().String(), splitContent[1]), nil
}

func MakeUrl(bucket string, key string) (*url.URL, error) {
  return url.ParseRequestURI(fmt.Sprintf("https://%s/%s", bucket, key))
}

func Storage(key string, data *bytes.Reader, contentType string) (*url.URL, error) {
  s, err := session.NewSession(&aws.Config{Region: aws.String(getValue(DEFAULT_REGION, OS_REGION))})
  if err != nil {return nil, err}
  bucket := getValue(DEFAULT_BUCKET, OS_BUCKET)

  _, err = s3.New(s).PutObject(&s3.PutObjectInput{
    Bucket: aws.String(bucket),
    Key: aws.String(key),
    Body: data,
    ContentType: aws.String(contentType),
  })
  if err != nil {
    return nil, err
  }
  return MakeUrl(bucket, key)
}
