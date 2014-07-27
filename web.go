package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "os"
)

func main() {
  http.HandleFunc("/", profiles)
  fmt.Println("listening...")
  err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
  if err != nil {
    panic(err)
  }
}

func profiles(w http.ResponseWriter, r *http.Request) {
  profile := Profile{"Alex", []string{"Cleiton Francisco", "Programador"}}

  js, err := json.Marshal(profile)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func hello(res http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(res, "hello, world")
}
