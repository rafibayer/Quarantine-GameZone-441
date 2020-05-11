# Quarantine-GameZone-441
Informatics 441 (Server-side) final project

# Project Proposal
https://docs.google.com/document/d/11USZUFFNMCr2rI1_OYwn5Vvh6M2UfSBn8hGCUdOAKZw/edit?usp=sharing

# Proposal Contents




# Project description
Quarantine Game-Zone is an application that allows you to play games with others online via a web browser.
Our target audience is anyone looking for a way to stay social by playing games, especially during quarantine. There are many other similar services, such as Jackbox, Drawful, and more, but we want to offer a free and easy-to-use alternative. As developers, we want to create an app that we could see ourselves using. As students experiencing this unprecedented online quarter, we think it’s more important than ever to stay connected, and games are a fantastic way to do that.


# Endpoints
Games lobby
   * /v1/games/
      * GET; Admin purposes, see all currently running games
         * 200: Gets all games sessions that are currently happening, returns list
         * 500: Internal server error
* /v1/games/tictactoe 
   * POST
      * 201 created: Creates a game state on the server, sends you the initial state of the game as JSON
      * 500: Internal server error
* /v1/games/tictactoe/{lobby id}
   * GET
      * 200 ok: Returns the current state of the game
      * 401 unauthorized: Could not verify player, or they are not in the game
      * 404 not found: The game wasn’t found
      * 415: Unsupported media type
      * 500: Internal server error
   * POST
      * 201 created: Applies the move to the game, returns the updated game state as JSON
      * 400 bad request: An illegal move is given
      * 401 unauthorized: Could not verify player, or they are not in the game 
      * 404 not found: The game wasn’t found
      * 415: Unsupported media type
      * 500: Internal server error
   * DELETE
      * 200: Successfully ends the game session
      * 401: Could not verify player or they are are not in the game
      * 404: The game wasn’t found
      * 415: Unsupported media type
Leaderboard
* /v1/leaderboard?game={GAME}&top={NUM}
   * GET
      * 200: returns top NUM (or default) players and number of wins for GAME
      * 400: game not found
      * 400: num invalid or too big
      * 500: Internal server error
* /v1/leaderboard/{playerid}?game={GAME}
   * GET
      * 200: returns the number of wins for player for a game (in context of playing a game)
      * 404: Player not found
Players
* /v1/players
   * POST
      * 201 created: Create a new player
Specific player
   * /v1/players/{player id OR me}
      * PATCH
         * 200 ok:  update player (first name, last name)
         * 403 forbidden: not authenticated to make changes to this player profile
         * 404 not found: player not found
         * 415: Unsupported media type
      * GET
         * 200 ok: get player info
         * 403 forbidden: not authenticated to get player profile
         * 404 not found: player not found
         * 415: Unsupported media type
      * DELETE
         * 200 ok:
         * 403 forbidden: not authenticated to delete this player
         * 404 not found: player not found
         * 415: Unsupported media type
Sessions
   * /v1/sessions
      * POST
         * 201 created: Created a new session
         * 401 unauthorized: Bad credentials
         * 415: unsupported media type
         * 500: Internal server error
Specific session
   * /v1/sessions/{session id or mine}
      * DELETE
         * 403 forbidden: not mine
      * End session


# Models


* Player
   * Player profile
      * Player ID
      * User Name
      * Email
      * First Name
      * Last Name
      * Password Hash
   * Player create account
      * User Name
      * Email
      * First Name
      * Last Name
      * Password
      * PasswordConf
   * Credentials
      * User Name
      * Password 
* Game
   * Tic-tac-toe
      * Lobby ID
      * X PlayerID
      * O PlayerID
      * Gamestate
         * Whose turn
         * Outcome
         * Board representation 
* Leaderboard
   * PlayerID
   * GameName
   * Score


# Use cases and priority:

Priority
	User
	Description
	P0
	As a player
	I want to be able to create a game and share it with friends
	P0
	As a player
	I want to be able to join a created game
	P0
	As a player
	I want to be able to play the game until I choose to leave or a winner is decided
	P1
	As a player
	I want to be able to view leaderboards for each game
	P2
	As a player
	I want to be able to chat with other players within a game
	





# Infrastructure diagram


https://app.lucidchart.com/invitations/accept/ffb7c05e-ab8e-4cce-aa82-9e2046c505b6