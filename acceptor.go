package main

import (
  "github.com/julienschmidt/httprouter"
  "net/http"
)

func main() {
  r := httprouter.New()
  
  r.GET("/_/*p", process)
  r.POST("/_/*p", process)
  
  http.ListenAndServe("0.0.0.0:8080", r)
}

