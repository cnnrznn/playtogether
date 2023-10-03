package api

import (
	"encoding/json"
	"net/http"

	"github.com/cnnrznn/playtogether/api/model"
)

func Run() error {
	http.HandleFunc("/play", HandlePlay)

	return http.ListenAndServe(":8080", nil)
}

func HandlePlay(w http.ResponseWriter, req *http.Request) {
	var ping model.Ping

	err := json.NewDecoder(req.Body).Decode(&ping)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeResponse(w, WebRes{
		Status:  http.StatusOK,
		Payload: ping,
	})
}

type WebErr struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type WebRes struct {
	Status  int `json:"status"`
	Payload any `json:"payload"`
}

func writeError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)

	bs, err := json.Marshal(WebErr{
		Status: code,
		Error:  err.Error(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"problem serializing response\"}"))
		return
	}

	w.Write(bs)
}

func writeResponse(w http.ResponseWriter, payload WebRes) {
	w.WriteHeader(http.StatusOK)

	bs, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":\"problem serializing response\"}"))
		return
	}

	w.Write(bs)
}
