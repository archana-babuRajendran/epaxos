package main

import (
  "log"
  "net/http"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Fielding request")
  w.Write([]byte("Youve reached state"))
}

func main() {
  mux := http.NewServeMux()
  mux.Handle("/state/1", rh)
  log.Println("Listening...")
  http.ListenAndServe(":5000", mux)
}
