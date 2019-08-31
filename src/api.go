package main

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"github.com/go-chi/chi"
)

type index struct {
	Stat int `json:"stat"`
	Word string `json:"word"`
}

func getIndexStat(w http.ResponseWriter, r *http.Request) {
	word := chi.URLParam(r, "word")

	stat, _ := StatsStore.GetStat(word)
	
	json.NewEncoder(w).Encode(index{Word: word, Stat: stat})
}

func inputHandler(w http.ResponseWriter, r *http.Request) {
	inputType := r.URL.Query().Get("source")
	var input string
	var err error

	if inputType == Text.String() {
		// Since body may be large we first upload it to temp file
		input, err = uploadBody(r.Body)
		inputType = File.String()
	} else {
		var data []byte
		data, err = ioutil.ReadAll(r.Body)
		input = string(data)
	}

	if err != nil {
		w.WriteHeader(400)
		return
	}

	switch strings.ToLower(inputType) {
	case URL.String():
		err = NewBatchingJob(string(input), URL).Dispatch()
	case File.String():
		err = NewBatchingJob(string(input), File).Dispatch()
	default:
		w.Write([]byte("Unsupported input type"))
		w.WriteHeader(400)
		return
	}

	if err != nil {
		w.WriteHeader(500)
	}
}

// StartServer creates new http listener on specific port
func StartServer(port int, broker Broker) {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})
	

	r.Route("/index", func(r chi.Router) {
		r.Post("/", inputHandler)
		r.Get("/{word}", getIndexStat)
	})
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}