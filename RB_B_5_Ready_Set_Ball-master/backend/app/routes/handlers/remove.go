package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// Remove handles requests to remove items from a user
// to remove a friend with username mich:
//		url - /remove/f
//		params - "result" : {"unfriend" : "mich"}
func (env *Env) Remove(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()
	log.Printf("Store is %#v", store)

	//session
	sess, err := env.sessionStore.Get(req, config.PostSessionName)
	if !sess.IsAuth() {
		yoda.SendClientError(w, yoda.ErrNotLoggedIn, http.StatusUnauthorized)
		return
	}

	params := retrieveParamsFromContext(req.Context())
	if params == nil {
		yoda.SendClientError(w, yoda.ClientError{Code: 1, Message: "Invalid parameters"}, 0) //TODO: use proper code
		return
	}
	yReq, err := yoda.NewRequestWithSocket(req, params, sess.All(), env.socketManager)
	if err != nil {
		yoda.SendClientError(w, err, 0) //TODO: add code
		return
	}

	result, err := operations.Do("remove", store, *yReq)
	if err != nil {
		sess.Save(req, w)               //TODO: why save this session?
		yoda.SendClientError(w, err, 0) //TODO: add code
		return
	}

	sess.Save(req, w)
	yRes := yoda.NewResponse(result.GetMeta(), result.GetData(), false) //TODO use one argument newresponse
	json.NewEncoder(w).Encode(yRes)
}
