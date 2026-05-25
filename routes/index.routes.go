package routes

import "net/http"

var helloGO string = "Hello Go, This is MY first API rest 40"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(helloGO))
}
