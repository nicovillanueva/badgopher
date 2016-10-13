package main

import (
	"fmt"
	//	"io/ioutil"
	"net/http"
)

func encryptHandler(w http.ResponseWriter, r *http.Request) {
	k := r.FormValue("key")
	p := r.FormValue("path")
	fmt.Println(k, p)
	//WalkPath(p, []byte(k), true)
}

func decryptHandler(w http.ResponseWriter, r *http.Request) {
	k := r.FormValue("key")
	p := r.FormValue("path")
	fmt.Println(k, p)
	//WalkPath(p, []byte(k), false)
}

func serve() {
	http.HandleFunc("/encrypt", encryptHandler)
	http.HandleFunc("/decrypt", decryptHandler)
	http.ListenAndServe(":8080", nil)
}
