package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
)

// GetUsers handles request to get all users
func (env *Env) GetUsers(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()
	log.Printf("Store is %#v", store)

	var yRes *yoda.Response
	params := retrieveParamsFromContext(req.Context()) //TODO: check for error
	yReq, err := yoda.NewRequest(req, params)

	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	result, err := operations.Do("get_users", store, *yReq)
	if err != nil {
		yoda.SendClientError(w, err, 0)
		return
	}

	yRes = yoda.NewResponse(result.GetMeta(), result.GetData(), false)
	json.NewEncoder(w).Encode(yRes)
}
