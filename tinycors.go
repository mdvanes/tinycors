package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
  "net/http"
  "reflect"
  "errors"
)

const (
	defaultPort = "3000"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

type arrayOrigins []string

func (i *arrayOrigins) String() string {
  return "my string representation"
}

func (i *arrayOrigins) Set(value string) error {
  *i = append(*i, value)
  return nil
}

var allowedOrigins arrayOrigins

func main() {
  var port = flag.String("port", defaultPort, "the port of the server")
  flag.Var(&allowedOrigins, "origins", "allowed origins, e.g. -origins http://localhost:3000 -origins http://localhost:8080")
  flag.Parse()
  fmt.Printf("origins: %+q\n", allowedOrigins)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)

    origin := r.Header.Get("origin")
    
    err := checkOrigin(w, origin)
    if err != nil {
			respondWithErr(w, err.Error())
      return
    }

		queryUrl := r.URL.Query().Get("url")
		if queryUrl == "" {
			respondWithErr(w, "empty url")
			return
		}

    resp, err := http.Get(queryUrl)

		if err != nil {
			respondWithErr(w, err.Error())
			return
		}
		defer func() {
			if resp.Body.Close() != nil {
				log.Println(err)
			}
		}()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			respondWithErr(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(body); err != nil {
			log.Println(err)
		}

	})

	log.Println("Starting TinyCORS server on", *port)
	if err := http.ListenAndServe(":" + *port, nil); err != nil {
		log.Fatal()
	}
}

func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

  // TODO
	// if arr.Kind() != reflect.Array {
	// 	panic("Invalid data-type")
	// }

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func enableCors(w *http.ResponseWriter) {
  (*w).Header().Set("Access-Control-Allow-Origin", "*")
  (*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

// If supplied originIncludeList is set use that, otherwise allow everything
func checkOrigin(w http.ResponseWriter, origin string) error {
  if len(allowedOrigins) > 0 && !itemExists(allowedOrigins, origin) {
    log.Println("origin", origin, "is not in includelist")
    w.WriteHeader(http.StatusForbidden)
    return errors.New("origin is not in includelist")
  }
  return nil
}

func respondWithErr(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println("some error", err)
	// _, _ = w.Write([]byte(fmt.Sprintf(fmt.Sprintf(`{"err":"%v"}`, err))) // TODO fix
}
