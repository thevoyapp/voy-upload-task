package main

import (
  "bytes"
  "net/http"
  "net/url"
  "io"
)

type Audio struct {
  Body  *bytes.Reader
}

func (a *Audio) Read() []byte {
  data := make([]byte, a.Body.Size())
  a.Body.Read(data)
  a.Body.Seek(0, io.SeekStart)
  return data
}

func NewAudio(data []byte) *Audio {
  return &Audio{
    Body: bytes.NewReader(data),
  }
}

func (a *Audio) ContentType() string {
  return http.DetectContentType(a.Read())
}

func (a *Audio) Upload(path string) (*url.URL, error) {
  contentType := a.ContentType()
  key, err := FileKey(path, contentType)
  if err != nil {
    return nil, err
  }
  return Storage(key, a.Body, contentType)
}
