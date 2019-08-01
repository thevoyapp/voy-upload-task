package main

import (
  "bytes"
  "gopkg.in/h2non/bimg.v1"
  "log"
  "net/http"
  "io"
  "net/url"
)

const (
  MAX_PIXEL = 1080
)

type Image struct {
  Body    *bytes.Reader
}

func max(n1 int, n2 int) int {
  if n1 > n2 {return n1}
  return n2
}

func (im *Image) Read() []byte {
  data := make([]byte, im.Body.Size())
  im.Body.Read(data)
  im.Body.Seek(0, io.SeekStart)
  return data
}

func NewImage(data []byte) *Image {
  return &Image{
    Body: bytes.NewReader(data),
  }
}

func (im *Image) Compress() {
  data := im.Read()

  reprImage := bimg.NewImage(data)
  size, _ := reprImage.Size()
  height := size.Height
  width := size.Width

  largest := max(height, width)
  if largest > MAX_PIXEL {
    height = height * MAX_PIXEL / largest
    width = width * MAX_PIXEL / largest
  }

  result, err := reprImage.Process(
    bimg.Options{
      Gravity: bimg.GravitySmart,
      Quality: 50,
      NoAutoRotate: true,
      Crop: true,
      Height: height,
      Width: width,
    })

  if err != nil {
    log.Print("Unable to process image: " + err.Error())
  } else {
    im.Body = bytes.NewReader(result)
  }
}

func (im *Image) ContentType() string {
  return http.DetectContentType(im.Read())
}

func (im *Image) Upload(path string) (*url.URL, error) {
  im.Compress()
  contentType := im.ContentType()
  key, err := FileKey(path, contentType)
  if err != nil {
    return nil, err
  }
  return Storage(key, im.Body, contentType)
}
