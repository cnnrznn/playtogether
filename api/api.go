package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cnnrznn/playtogether/model"
	"github.com/cnnrznn/playtogether/service/game"
	"github.com/cnnrznn/playtogether/service/play"
	"github.com/google/uuid"
)

type WebErr struct {
	Error string `json:"error"`
}

type WebRes struct {
	Payload any `json:"payload"`
}

func Run() error {
	http.HandleFunc("/play", HandlePlayRequest)
	http.HandleFunc("/game", HandleGame)

	return http.ListenAndServe(":8080", nil)
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

func HandleGame(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	gr := GameRequest{}

	if err := json.NewDecoder(req.Body).Decode(&gr); err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	switch gr.Action {
	case "create":
		CreateGame(w, req, gr)
	case "confirm":
		ConfirmGame(w, req, gr)
	case "get":
	default:
		writeError(w, fmt.Errorf("unsupported 'game' action"), http.StatusBadRequest)
		return
	}
}

func ConfirmGame(w http.ResponseWriter, req *http.Request, gr GameRequest) {
	if err := game.Confirm(gr.GameID, gr.PlayRequestID); err != nil {
		writeError(w, fmt.Errorf("could not confirm game: %w", err), http.StatusInternalServerError)
		return
	}

	writeResponse(w, WebRes{})
}

func CreateGame(w http.ResponseWriter, req *http.Request, gr GameRequest) {
	prIDs := make(map[uuid.UUID]struct{})

	for _, id := range gr.PlayRequestIDs {
		prIDs[id] = struct{}{}
	}

	if err := game.Create(model.Game{
		PlayRequests: prIDs,
	}); err != nil {
		writeError(w, fmt.Errorf("could not create game: %w", err), http.StatusInternalServerError)
		return
	}

	writeResponse(w, WebRes{})
}

func HandlePlayRequest(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch req.Method {
	case "GET":
		GetPlayRequest(w, req)
	case "POST":
		PostPlayRequest(w, req)
	}
}

func GetPlayRequest(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("user")
	if len(userID) == 0 {
		writeError(w, fmt.Errorf("must supply arg 'user'"), http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		writeError(w, fmt.Errorf("could not parse 'user' as UUID"), http.StatusInternalServerError)
		return
	}

	activity := req.URL.Query().Get("activity")
	if len(activity) == 0 {
		writeError(w, fmt.Errorf("must supply arg 'activity'"), http.StatusBadRequest)
		return
	}

	prs, err := play.GetPlayRequests(uid, activity)
	if err != nil {
		writeError(w, fmt.Errorf("error loading play requests: %w", err), http.StatusInternalServerError)
		return
	}

	writeResponse(w, WebRes{Payload: prs})
}

func PostPlayRequest(w http.ResponseWriter, req *http.Request) {
	var pr model.PlayRequest

	if err := json.NewDecoder(req.Body).Decode(&pr); err != nil {
		writeError(w, fmt.Errorf("could not decode json: %w", err), http.StatusInternalServerError)
		return
	}

	if err := play.CreatePlayRequest(pr); err != nil {
		writeError(w, fmt.Errorf("service err: %w", err), http.StatusInternalServerError)
		return
	}

	writeResponse(w, WebRes{})
}

/*
func HandleGames(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
*/
