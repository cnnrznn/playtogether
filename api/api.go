package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cnnrznn/playtogether/model"
	"github.com/cnnrznn/playtogether/service"
	"github.com/google/uuid"
)

func Run() error {
	http.HandleFunc("/ping", HandlePing)
	http.HandleFunc("/games", HandleGames)

	return http.ListenAndServe(":8080", nil)
}

func HandleGames(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		writeError(w, fmt.Errorf("bad method for /games endpoint"), http.StatusBadRequest)
		return
	}

	var player model.Player

	id := req.URL.Query().Get("id")
	if len(id) == 0 {
		writeError(w, fmt.Errorf("missing 'id' parameter"), http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		writeError(w, fmt.Errorf("couldn't parse uuid"), http.StatusInternalServerError)
		return
	}

	player.ID = uid

	games := service.GetPlayerGames(player)

	writeResponse(w, WebRes{
		Payload: games,
	})
}

func HandlePing(w http.ResponseWriter, req *http.Request) {
	var ping model.Ping

	err := json.NewDecoder(req.Body).Decode(&ping)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response, err := service.Update(ping)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeResponse(w, WebRes{
		Payload: response,
	})
}

type WebErr struct {
	Error string `json:"error"`
}

type WebRes struct {
	Payload any `json:"payload"`
}

func writeError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)

	bs, err := json.Marshal(WebErr{
		Error: err.Error(),
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
