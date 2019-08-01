package main

import (
  "net/http"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ssm"
  "log"
  "os"
)

const (
  OS_VERSION string = "VERSION"
  PARAM_PATH string = "VERSION_PARAMETER_PATH"
  TEMP_PARAM_PATH string = "TEMP_VERSION_PARAMETER_PATH"
)

func GetEnviron(environ string) (string, bool) {
  path, ok := os.LookupEnv(environ)
  if !ok {
    log.Printf("Parameter %s not found", environ)
    return "", false
  }
  return path, true
}

func GetParam(environ string) (string, bool) {
  svc := ssm.New(session.New())
  getInput := ssm.GetParameterInput{}
  path, ok := GetEnviron(environ)
  if !ok {return "", false}
  getInput.SetName(path)
  getInput.SetWithDecryption(false)
  parameter, err := svc.GetParameter(&getInput)
  if err != nil {
    log.Print(err.Error())
    return "", false
  }
  return *parameter.Parameter.Value, true
}

func HealthCheck(w http.ResponseWriter, r *http.Request)  {
  baseVersion, ok := GetParam(PARAM_PATH)
  if !ok {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  myVersion, ok := GetEnviron(OS_VERSION)
  if !ok {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  if baseVersion != myVersion {
    newVersion, ok := GetParam(TEMP_PARAM_PATH)
    if !ok {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    if myVersion != newVersion {
      log.Printf("Version %s is not %s or %s.", myVersion, baseVersion, newVersion)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
  }
  w.WriteHeader(http.StatusOK)
}
