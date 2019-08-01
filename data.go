package main

import (
  "bytes"
  "net/http"
  "strings"
  "errors"
  "net/url"
)

const (
  IMAGE_PREFIX string = "image"
  AUDIO_PREFIX string = "audio"
)

type DataType interface {
  Upload(string) (*url.URL, error)
}

func NewUploadData(data *bytes.Buffer) (DataType, error) {
  byteData := data.Bytes()
  mimeType := http.DetectContentType(byteData)
  if strings.HasPrefix(mimeType, IMAGE_PREFIX) {
    return NewImage(byteData), nil
  }
  if strings.HasPrefix(mimeType, AUDIO_PREFIX) {
    return NewAudio(byteData), nil
  }
  return nil, errors.New("No such data type: " + mimeType)
}
