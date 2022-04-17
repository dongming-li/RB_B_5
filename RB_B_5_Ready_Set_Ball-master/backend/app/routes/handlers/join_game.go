package handlers

import (
	"encoding/json"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// JoinGame handles request to join a game based on a joincode
// 		url - /game/join/mode
// made should be an i for ID or a j for Join Code
//params "code" = joinCode or gameID
// The result is a game
//		"result": game
func (env *Env) JoinGame(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()

	sess, _ := env.sessionStore.Get(req, config.PostSessionName)
	if !sess.IsAuth() {
		yoda.SendClientError(w, yoda.ErrNotLoggedIn, http.StatusUnauthorized)
		sess.Save(req, w)
		return
	}

	var yRes *yoda.Response
	params := retrieveParamsFromContext(req.Context())
	yReq, err := yoda.NewRequestWithSocket(req, params, sess.All(), env.socketManager)

	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	result, err := operations.Do("join_game", store, *yReq)
	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}
	sess.Save(req, w)
	yRes = yoda.NewResponse(result.GetMeta(), result.GetData(), false)
	json.NewEncoder(w).Encode(yRes)
}
