package api

import "net/http"

func Run() error {
	http.HandleFunc("/play", HandlePlay)

	return http.ListenAndServe(":8080", nil)
}

func HandlePlay(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
