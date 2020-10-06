package main

import (
	"flag"
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	address = ":9009"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func main() {
	var addr = flag.String("addr", address, "the address of the application")
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

	log.Println("Starting web server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal()
	}
}

func respondWithErr(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println("some error", err)
	// _, _ = w.Write([]byte(fmt.Sprintf(fmt.Sprintf(`{"err":"%v"}`, err))) // TODO fix
}
