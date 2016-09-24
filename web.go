package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getDaKey() []byte {
	k, _ := ioutil.ReadFile("key.log")
	return k
}

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

func launch() {
	http.HandleFunc("/encrypt", encryptHandler)
	http.HandleFunc("/decrypt", decryptHandler)
	http.ListenAndServe(":8080", nil)
}
