package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "os"
  "log"
  "net/http"
  "strconv"

  "github.com/mattbaird/elastigo/api"
  "github.com/mattbaird/elastigo/core"
  "github.com/gorilla/mux"
)

// error response contains everything we need to use http.Error
type handlerError struct {
  Error   error
  Message string
  Code    int
}

// Book model
type Book struct {
  ID     string `json:"id"`
  Title  string `json:"title"`
  Author string `json:"author"`
  Image  string `json:"image"`
}

// Result
type Result struct {
  Total int    `json:"total"`
  Page  int    `json:"page"`
  Query string `json:"query"`
  Books []Book `json:"books"`
}

type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // call the actual handler
  response, err := fn(w, r)

  // check for errors
  if err != nil {
    log.Printf("ERROR: %v\n", err.Error)
    http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
    return
  }
  if response == nil {
    log.Printf("ERROR: response from method is nil\n")
    http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
    return
  }

  // turn the response into JSON
  bytes, e := json.Marshal(response)
  if e != nil {
    http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
    return
  }

  // send the response and log
  w.Header().Set("Content-Type", "application/json")
  w.Write(bytes)
  log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}

func config() {
  api.Protocol = "https"
  api.Domain   = "pepper-8677672.us-east-1.bonsai.io"
  api.Username = "8oxyrtww"
  api.Password = "h7r6vjg2648whdmb"
  api.Port     = "443"
}

func mountJSON(query string, page int) string {
  var json string
  size := 4
  from := size * (page - 1)

  if query == "" {
    json = fmt.Sprintf(`{
      "from" : %d, 
      "size" : %d,
      "query": {
        "match_all": {}
      }
    }`, from, size)
  } else {
    json = fmt.Sprintf(`{
      "from" : %d, 
      "size" : %d,
      "query": {
        "simple_query_string" : {
          "query":  "%s",
          "fields": [ "title", "author" ],
          "default_operator": "and"
        }
      }
    }`, from, size, query)
  }
  return json
}

func searchBooks(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
  result := Result{Total: 0, Page: 1, Query: "", Books: make([]Book, 0)}
  params := r.URL.Query()
  result.Query = params.Get("q")

  page, err_page := strconv.Atoi(params.Get("p"))
  if (err_page != nil) {
    page = 1
  }

  searchJSON := mountJSON(result.Query, page)

  out, e := core.SearchRequest("books", "book", nil, searchJSON)
  if e != nil {
    return result, &handlerError{e, "JSON error", http.StatusBadRequest}
  }

  result.Page  = page
  result.Total = out.Hits.Total

  for _, hit := range out.Hits.Hits {
    var b Book
    if err := json.Unmarshal(*hit.Source, &b); err == nil {
      result.Books = append(result.Books, b)
    }
  }

  return result, nil
}

func main() {
  dir := flag.String("directory", "web/", "directory of web files")
  flag.Parse()

  // handle all requests by serving a file of the same name
  fs := http.Dir(*dir)
  fileHandler := http.FileServer(fs)
  config()

  // setup routes
  router := mux.NewRouter()
  router.Handle("/", http.RedirectHandler("/bookstore/", 302))
  router.Handle("/books", handler(searchBooks)).Methods("GET")
  router.PathPrefix("/bookstore/").Handler(http.StripPrefix("/bookstore", fileHandler))
  http.Handle("/", router)

  log.Println("Running...")

  err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
  fmt.Println(err.Error())
}
