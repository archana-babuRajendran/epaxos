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
  mux.HandleFunc("/state/1", dummyHandler)
  log.Println("Listening...")
  log.Fatal(http.ListenAndServe(":5000", mux))
}
