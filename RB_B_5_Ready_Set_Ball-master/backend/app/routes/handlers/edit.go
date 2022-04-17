package handlers

import (
	"encoding/json"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// Edit handles requests to modify properties of an entity
// to edit the user properties (username and email) of the logged in user:
//		url - /edit/user
//		params - "result" : {"username" : "mich", "email": "abc@de.com"}
//		url - /edit/game
// 		params - {name: newName, minAge: 15, maxAge: 24, lat: 15, lng: 24, sport:3}
// If you don't want that thing changed don't inculde it when sending
func (env *Env) Edit(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()

	//session
	sess, err := env.sessionStore.Get(req, config.PostSessionName)
	if !sess.IsAuth() {
		yoda.SendClientError(w, yoda.ErrNotLoggedIn, http.StatusUnauthorized)
		return
	}

	params := retrieveParamsFromContext(req.Context())
	if params == nil {
		yoda.SendClientError(w, yoda.ErrInvalidParameters, http.StatusBadRequest)
		return
	}
	yReq, err := yoda.NewRequestWithSocket(req, params, sess.All(), env.socketManager)
	if err != nil {
		yoda.SendClientError(w, err, 0) //TODO: add code
		return
	}

	result, err := operations.Do("edit", store, *yReq)
	if err != nil {
		sess.Save(req, w) //TODO: why save this session?
		yoda.SendClientError(w, err, 0)
		return
	}

	// update session with new user values
	if params["entity"] == "user" {
		sess.Assign(result)
	}
	sess.Save(req, w)
	yRes := yoda.NewResponse(result.GetMeta(), result.GetData(), false) //TODO use one argument newresponse
	json.NewEncoder(w).Encode(yRes)
}
