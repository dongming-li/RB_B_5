package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// ExitGame handles request to exit a game based on user's currentgame
// 		url - /game/exit
// The result is a message saying it succeeded
func (env *Env) ExitGame(w http.ResponseWriter, req *http.Request) {
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
	yReq, err := yoda.NewRequestWithSocket(req, params, sess.All(), env.socketManager)

	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	result, err := operations.Do("exit_game", store, *yReq)
	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}
	sess.Save(req, w)
	yRes = yoda.NewResponse(result.GetMeta(), result.GetData(), false)
	json.NewEncoder(w).Encode(yRes)
}
