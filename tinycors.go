package main

import (
	"flag"
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	defaultPort = "3000"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func main() {
	var port = flag.String("port", defaultPort, "the port of the server")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

func respondWithErr(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println("some error", err)
	// _, _ = w.Write([]byte(fmt.Sprintf(fmt.Sprintf(`{"err":"%v"}`, err))) // TODO fix
}
