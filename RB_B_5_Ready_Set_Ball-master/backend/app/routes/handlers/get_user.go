package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// GetUser handles request to get a person by username or id
// Route: "/user/:command/:argument"
// when command is f:
//      all friends of the user is returned
//			url - /user/f/#
//			params - "result" : {"username" : "mich"} => returns mich's friends
// when command is pg:
//			url - /user/pg/#
//			params - "result" : {"username" : "mich"} => returns mich's prevgames
// when command is p:
//      populate = 0, the user is returned without friends
//			url - /user/p/0
//			params - "result" : {"username" : "mich"} => returns mich user object
//		populate = 1, the user is returned with friend
//			url - /user/p/1
//			params - "result" : {"username" : "mich"} => returns mich user object with friends
// when there is no username in the body, the signed in user is used
func (env *Env) GetUser(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()
	log.Printf("Store is %#v", store)

	sess, err := env.sessionStore.Get(req, config.PostSessionName)
	if !sess.IsAuth() {
		yoda.SendClientError(w, yoda.ErrNotLoggedIn, http.StatusUnauthorized)
		return
	}

	var yRes *yoda.Response

	params := retrieveParamsFromContext(req.Context())
	if params == nil {
		yoda.SendClientError(w, yoda.ClientError{Code: 1, Message: "Invalid parameters"}, http.StatusBadRequest)
		return
	}

	yReq, err := yoda.NewRequestWithSession(req, params, sess.All())
	if err != nil {
		yoda.SendClientError(w, err, 0) //TODO: add code
		return
	}

	result, err := operations.Do("friend_user", store, *yReq)
	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	sess.Save(req, w)
	yRes = yoda.NewResponse(result.GetMeta(), result.GetData(), false)
	json.NewEncoder(w).Encode(yRes)
}
