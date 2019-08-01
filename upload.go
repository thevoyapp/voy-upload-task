package main

import (
  "net/http"
  "net/url"
  "bytes"
  "io"
  "log"
  "encoding/json"
  voymiddleware "github.com/CodyPerakslis/voy-middleware"
)

type TotalResponses struct {
  Results    []*FileResults   `json:"files"`
}

type FileResults struct {
  Filename   string      `json:"filename"`
  Url        *url.URL     `json:"url"`
}

func UploadContent(w http.ResponseWriter, r *http.Request)  {
  reader, err := r.MultipartReader()
  if err != nil {
    log.Print(err.Error())
    http.Error(w, "Could not read input", http.StatusInternalServerError)
    return
  }
  user, err := voymiddleware.GetUser(r)
  if err != nil {
    log.Print(err.Error())
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
  }

  path, err := user.GetAccessPath()
  if err != nil {
    log.Print(err.Error())
    http.Error(w, "No User Access Path: " + err.Error(), http.StatusUnauthorized)
    return
  }

  var buf bytes.Buffer
  count := 0
  responses := make(chan FileResults)
  for {
    part, err := reader.NextPart()
    if err == io.EOF {break}

    io.Copy(&buf, part)
    data, err := NewUploadData(&buf)
    buf.Reset()
    if err != nil {
      log.Print("Incorrect Format: " + err.Error())
      continue
    }

    go func() {
      result := FileResults{
        Filename: part.FileName(),
        Url: nil,
      }
      url, err := data.Upload(*path)
      if err != nil {
        log.Print("Data Upload Error: " + err.Error())
      } else {
        result.Url = url
      }
      responses <- result
    }()
    count++
  }

  totalResponses := []*FileResults{}
  for i := 0; i < count; i++ {
    result := <- responses
    log.Print(result.Filename, result.Url)
    totalResponses = append(totalResponses, &result)
  }

  responseBody, err := json.Marshal(
    TotalResponses{
      Results: totalResponses,
    })
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Content-Type", "application/json")
  w.Write(responseBody)
}
