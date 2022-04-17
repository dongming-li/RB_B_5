package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// GetGamesNearLocation handles request to get a games based on a location
// to get games around a region of (lng: 15, lat: 10),
// 		url - /games/l/lng/15/lat/10
// The result is a slice of games
//		"result": [gameA, gameB, ..]
func (env *Env) GetGamesNearLocation(w http.ResponseWriter, req *http.Request) {
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

	result, err := operations.Do("get_games_from_location", store, *yReq)
	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	yRes = yoda.NewResponse(result.GetMeta(), result.GetData(), false)
	json.NewEncoder(w).Encode(yRes)
}
