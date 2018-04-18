package main

//go:generate go-bindata -o ./bindata.go data/...
import (
  "net/http"
  "strings"
  "fmt"
)

func respond(w http.ResponseWriter, r *http.Request) {
  message := r.URL.Path
  message = strings.TrimPrefix(message, "/")
  message = "Hello " + message
  w.Write([]byte(message))
}

func main() {
	data, err := Asset("data/foo.css")
	if err != nil {
			// asset was not found.
			fmt.Println("asset not found")
		}
	fmt.Println(data)

		// use asset data

  http.HandleFunc("/", respond)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}
