package main

import (
	"errors"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

const (
	defaultPort = "3000"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

type arrayOrigins []string

// TODO needed?
func (i *arrayOrigins) String() string {
	return "my string representation"
}

// TODO needed?
func (i *arrayOrigins) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var allowedOrigins arrayOrigins

func main() {
	var port = flag.String("port", defaultPort, "the port of the server")
	flag.Var(&allowedOrigins, "origins", "allowed origins, e.g. -origins http://localhost:3000 -origins http://localhost:8080")
	flag.Parse()
	log.Printf("Allowed origins: %+q\n", allowedOrigins)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		origin := r.Header.Get("origin")

		if len(r.URL.Query()) == 0 {
			log.Println("FOOBAR")
			respondWithErr(w, "Hier uw documentatie")
			return
		}

		err := checkOrigin(w, origin)
		if err != nil {
			respondWithErr(w, err.Error())
			return
		}

		queryUrl := r.URL.Query().Get("get")
		if queryUrl == "" {
			respondWithErr(w, "query param \"get\" is not set in "+r.URL.EscapedPath())
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

	log.Println("Starting TinyCORS 🌱 server on", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal("Failed to start TinyCORS server, port in use?")
	}
}

func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

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
		msg := "origin \"" + origin + "\" is not in includelist"
		// TODO should return 403 in respondWIthErr = w.WriteHeader(http.StatusForbidden)
		return errors.New(msg)
	}
	return nil
}

func respondWithErr(w http.ResponseWriter, err string) {
	log.Println("Error:", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	// TODO escape " to something readable, instead of html entity
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, html.EscapeString(err))))
}
