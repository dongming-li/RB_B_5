package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// Invites handles requests from one [model.User] sending and reviewing
// the request of/to another [model.User]
// to send friendRequest to a user with username mich:
//		url - /invite/m/send/t/0
//		params - "result" : {"to" : "mich"}
// to accept a friendRequest from a user with username mich:
//		url - /invite/m/review/t/0
//		params - "result" : {"from" : "mich", "accept": true}
// to reject a friendRequest from a user with username mich:
//		url - /invite/m/review/t/0
//		params - "result" : {"from" : "mich", "accept": false}
// to send a gameRequest to a user with username mich:
//		url- /invite/m/send/t/1
//		params - "result" : {"to" : "mich"}
// to accept a gameRequest from a user with username mich:
//		url - /invite/m/review/t/1
//		params - "result" : {"accept": true, "game": gameid}
// to reject a gameRequest from a user with username mich:
//		url - /invite/m/review/t/1
//		params - "result" : {"accept": false, "game": gameid}
func (env *Env) Invites(w http.ResponseWriter, req *http.Request) {
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

	result, err := operations.Do("invite", store, *yReq)
	if err != nil {
		sess.Save(req, w)               //TODO: why save this session?
		yoda.SendClientError(w, err, 0) //TODO: add code
		return
	}

	sess.Save(req, w)
	yRes := yoda.NewResponse(result.GetMeta(), result.GetData(), false) //TODO use one argument newresponse
	json.NewEncoder(w).Encode(yRes)
}
