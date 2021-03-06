package handlers

import (
	"Quarantine-GameZone-441/servers/gateway/gamesessions"
	"Quarantine-GameZone-441/servers/gateway/sessions"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// type GetAllResponse struct {
// 	map [GameLobby]map
// }

//LobbyHandlerPost handles request for making a game lobby
func (ctx *HandlerContext) LobbyHandlerPost(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "415: Request body must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	// decode request body into a new gamelobby
	decoder := json.NewDecoder(r.Body)
	newGameLobby := NewGameLobby{}
	err := decoder.Decode(&newGameLobby)
	if err != nil {
		http.Error(w, "Please provide a valid game lobby", http.StatusBadRequest)
		return
	}

	SessionState := SessionState{}
	playerSessID, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &SessionState)
	if err != nil {
		http.Error(w, "Please create a nickname to start your playing experience", http.StatusUnauthorized)
		return
	}
	playersSlice := make([]sessions.SessionID, 0)
	playersSlice = append(playersSlice, playerSessID)

	//create official game lobby
	gameLobby := &GameLobby{}
	if gameCapacity, Ok := GameCapacity[newGameLobby.GameType]; Ok {
		gameLobby.GameType = newGameLobby.GameType
		gameLobby.Capacity = gameCapacity
		gameLobby.Players = playersSlice
	} else {
		http.Error(w, fmt.Sprintf("we currently dont support the game: %s", newGameLobby.GameType), http.StatusBadRequest)
		return
	}
	//begins a session
	GameLobbyState := GameLobbyState{time.Now(), gameLobby}
	newGameSessionID, err := gamesessions.BeginGameSession(ctx.SigningKey, ctx.GameSessionStore, GameLobbyState, w)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//TODO: this creates a bug because the new gamesessionID is never stored in redis properly
	GameLobbyState.GameLobby.ID = newGameSessionID
	_, err = gamesessions.UpdateGameSession(ctx.SigningKey, ctx.GameSessionStore, GameLobbyState, w, newGameSessionID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ResponseGameLobby, err := ctx.convertToResponseLobbyForClient(*GameLobbyState.GameLobby)
	if err != nil {
		http.Error(w, "Please make sure all game players have a nickname", http.StatusUnauthorized)
		return
	}

	//Responds back to the user with the updated user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(ResponseGameLobby)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	return
}

// LobbyHandlerGet returns all public game lobbies
func (ctx *HandlerContext) LobbyHandlerGet(w http.ResponseWriter, r *http.Request) {
	SessionState := SessionState{}
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &SessionState)
	if err != nil {
		http.Error(w, "Please create a nickname", http.StatusUnauthorized)
		return
	}

	res, err := gamesessions.GetAllSessions(ctx.SigningKey, ctx.GameSessionStore)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error retrieving game lobbies", http.StatusInternalServerError)
		return
	}
	resultLobbies := make([]ResponseGameLobby, 0)
	for _, element := range res {
		gameLobby := GameLobby{}
		res, err := json.Marshal(element)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error retrieving game lobbies", http.StatusInternalServerError)
		}
		err = json.Unmarshal(res, &gameLobby)
		if len(gameLobby.GameID) == 0 {
			lobby, err := ctx.convertToResponseLobbyForClient(gameLobby)
			if err != nil {
				http.Error(w, "Error retrieving game lobbies", http.StatusInternalServerError)
				return
			}
			resultLobbies = append(resultLobbies, *lobby)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(resultLobbies)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
