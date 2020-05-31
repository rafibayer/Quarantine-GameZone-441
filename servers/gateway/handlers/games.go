package handlers

import (
	"Quarantine-GameZone-441/servers/gateway/gamesessions"
	"Quarantine-GameZone-441/servers/gateway/sessions"

	"net/http"
)

//TODO: Make sure we never send back user tokens or any other sensitive info
//	change it all to plain text repsonses later

//GameLobby struct represents the state of a game lobby, this is created for every game
type GameLobby struct {
	ID       gamesessions.GameSessionID `json:"game_id"`
	GameType string                     `json:"game_type"`
	Private  bool                       `json:"private"`
	Players  []sessions.SessionID       `json:"players"`
	Capacity int                        `json:"capacity"`
	GameID   string                     `json:"gameID"`
}

//NewGameLobby struct represents the state of a game lobby, this is created for every game
type NewGameLobby struct {
	GameType string `json:"game_type"`
	Private  bool   `json:"private"`
}

//ResponseGameLobby struct represents the state of the lobby that is sent to the client, with no session IDs
// instead the usern nicknames are stored
type ResponseGameLobby struct {
	ID       gamesessions.GameSessionID `json:"lobby_id"`
	GameType string                     `json:"game_type"`
	Private  bool                       `json:"private"`
	Players  []string                   `json:"players"`
	Capacity int                        `json:"capacity"`
	GameID   string                     `json:"gameID"`
}

// replaces all sessionIDs with player nicknames for client
func (ctx *HandlerContext) convertToResponseLobbyForClient(gameLobby GameLobby) (*ResponseGameLobby, error) {
	nicknames := make([]string, 0)
	gameLobbyResponse := &ResponseGameLobby{}
	for _, player := range gameLobby.Players {
		playerSessionState := &SessionState{}
		err := ctx.SessionStore.Get(player, playerSessionState)
		if err != nil {
			return nil, err
		}
		nicknames = append(nicknames, playerSessionState.Nickname)
	}
	gameLobbyResponse.ID = gameLobby.ID
	gameLobbyResponse.GameType = gameLobby.GameType
	gameLobbyResponse.Private = gameLobby.Private
	gameLobbyResponse.Capacity = gameLobby.Capacity
	gameLobbyResponse.GameID = gameLobby.GameID
	gameLobbyResponse.Players = nicknames
	return gameLobbyResponse, nil
}

func (ctx *HandlerContext) LobbyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		ctx.LobbyHandlerPost(w, r)
		return
	}
	if r.Method == http.MethodGet {
		ctx.LobbyHandlerGet(w, r)
		return
	}
	http.Error(w, "405: Method not allowed", http.StatusMethodNotAllowed)

}

//SpecificLobbyHandler handles request for a specific game session, currently we are supporting
// a post request that adds a new player to the game session, this is currently done by sending the sessID
// in the request body (maybe needs to change to getting the sessionID from autherization header)
func (ctx *HandlerContext) SpecificLobbyHandler(w http.ResponseWriter, r *http.Request) {

	GameSessionState := GameLobbyState{}
	_, err := gamesessions.GetGameState(r, ctx.SigningKey, ctx.GameSessionStore, &GameSessionState)
	if err != nil {
		http.Error(w, "game session doesn't exist", http.StatusUnauthorized)
		return
	}

	//this is meant to add a new player to the gamesession
	if r.Method == http.MethodPost {
		ctx.SpecificLobbyHandlerPost(w, r)
		return

	} else if r.Method == http.MethodGet {
		ctx.SpecificLobbyHandlerGet(w, r)
		return
	}
	http.Error(w, "405: Method not allowed", http.StatusMethodNotAllowed)
}

// /v1/gamelobby/lobbyid
// post
// all lobby changes (addingplayer)	-> creates a game (lets client know) ->
// -> client now knows to send get specific game /v1/game/lobbyid(Get) (start loop)
// get
// gets the specific lobby state
// patch
// removes players -> if 0 players, lobby deletes

// /v1/game/lobbyid
// post
//all gamestate changes (making a move)
// get
//gets the specific game state

// /v1/gamelobby
// post
// makes a lobby
// get
// gets all public games (lobby states)
