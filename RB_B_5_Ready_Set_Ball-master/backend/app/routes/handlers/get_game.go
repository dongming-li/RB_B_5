package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// GetGame handles request to get a game based on an ID
// 		url - /game/g/ID
// The result is a game
//		"result": game
func (env *Env) GetGame(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()
	log.Printf("Store is %#v", store)

	sess, _ := env.sessionStore.Get(req, config.PostSessionName)
	if !sess.IsAuth() {
		yoda.SendClientError(w, yoda.ErrNotLoggedIn, http.StatusUnauthorized)
		sess.Save(req, w)
		return
	}

	var yRes *yoda.Response
	params := retrieveParamsFromContext(req.Context()) //TODO: check for error
	yReq, err := yoda.NewRequestWithSession(req, params, sess.All())

	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	result, err := operations.Do("get_game", store, *yReq)
	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}
	sess.Save(req, w)
	yRes = yoda.NewResponse(result.GetMeta(), result.GetData(), false)
	json.NewEncoder(w).Encode(yRes)
}
