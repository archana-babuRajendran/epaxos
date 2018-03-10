package main

import (
  "log"
  "net/http",
  "client",
  "os/exec"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Fielding request")
  cmnd := exec.Command("../client/client", '-maddr "10.0.0.10" -mport 7087 -q 100 -c 100 -e true')
  cmnd.Start()
  w.Write([]byte("Youve reached state"))
}

func main() {
  mux := http.NewServeMux()
  mux.Handle("/state/1", rh)
  log.Println("Listening...")
  http.ListenAndServe(":5000", mux)
}
