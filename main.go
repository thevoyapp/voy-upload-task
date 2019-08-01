package main

import (
  "github.com/gorilla/mux"
  "net/http"
  voymiddleware "github.com/CodyPerakslis/voy-middleware"
)

func main() {
  base := mux.NewRouter()
  b := base.PathPrefix("/alive").Subrouter()
  b.HandleFunc("/", HealthCheck).
    Methods("GET").
    Name("Healthcheck")

  r := base.PathPrefix("/content").Subrouter()
  voymiddleware.AddApiMiddleware(r, true)
  voymiddleware.AddUserAuthorizedMiddleware(r, true)

  r.HandleFunc("/", UploadContent).
    Methods("POST").
    Name("Upload")

  http.ListenAndServe(":80", base)
}
