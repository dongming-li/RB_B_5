package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// UploadAvatar handles updating a user profile with a profile pic
// 		url - /upload
// The result is a success bool
//		"result": {"successful" : true}
func (env *Env) UploadAvatar(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()
	log.Printf("Store is %#v", store)

	sess, _ := env.sessionStore.Get(req, config.PostSessionName)
	if !sess.IsAuth() {
		yoda.SendClientError(w, yoda.ErrNotLoggedIn, http.StatusUnauthorized)
		sess.Save(req, w)
		return
	}

	params := retrieveParamsFromContext(req.Context())

	yReq, err := yoda.NewRequestFormData(req, params, sess.All())
	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	result, err := operations.Do("upload", store, *yReq)
	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	sess.Save(req, w)
	yRes := yoda.NewResponse(result.GetMeta(), result.GetData(), false)
	json.NewEncoder(w).Encode(yRes)
}
