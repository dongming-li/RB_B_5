package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// RateGame handles rating a game that has been played by a user
// 		url - /game/rate
//		To rate a game with ID, gameID a 3.5, the params will be:
//		params - {code: gameID, rating: 3.5}
// The result is a message saying it succeeded
func (env *Env) RateGame(w http.ResponseWriter, req *http.Request) {
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
	params := retrieveParamsFromContext(req.Context())
	yReq, err := yoda.NewRequestWithSession(req, params, sess.All())

	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	result, err := operations.Do("rate_game", store, *yReq)
	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}
	sess.Save(req, w)
	yRes = yoda.NewResponse(result.GetMeta(), result.GetData(), false)
	json.NewEncoder(w).Encode(yRes)
}
