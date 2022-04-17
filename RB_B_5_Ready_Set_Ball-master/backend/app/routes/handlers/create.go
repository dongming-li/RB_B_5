package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// CreateEntity handles request to create entities
func (env *Env) CreateEntity(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()
	log.Printf("Store is %#v", store)

	//session
	sess, _ := env.sessionStore.Get(req, config.PostSessionName)

	var yRes *yoda.Response
	params := retrieveParamsFromContext(req.Context()) //TODO: check for error
	if params == nil {
		yoda.SendClientError(w, yoda.ClientError{Code: 1, Message: "Invalid parameters"}, 0) //TODO: use proper code
		return
	}

	yReq, err := yoda.NewRequestWithSession(req, params, sess.All())
	if err != nil {
		yoda.SendClientError(w, err, 0) //TODO: add code
		return
	}

	if yReq.Params["entity"] != "user" && !sess.IsAuth() {
		yoda.SendClientError(w, yoda.ErrNotLoggedIn, http.StatusUnauthorized)
		sess.Save(req, w)
		return
	}
	result, err := operations.Do("create", store, *yReq)
	if err != nil {
		yoda.SendClientError(w, err, 0) //TODO: add code
		return
	}

	sess.Save(req, w)
	yRes = yoda.NewResponse(result.GetMeta(), result.GetData(), false) //TODO use one argument newresponse
	json.NewEncoder(w).Encode(yRes)
}
