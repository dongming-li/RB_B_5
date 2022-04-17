package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// Status handles request to send the status of a user (whether or not they are authenticated)
// Route: "/status"
//		returns {username: 'mich', isAuth: true} if authenticated
//		returns {isAuth: false} if not authenticated
func (env *Env) Status(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()
	log.Printf("Store is %#v", store)

	sess, _ := env.sessionStore.Get(req, config.PostSessionName)

	username, _ := sess.Get("username")
	result := yoda.NewResponse(map[string]string{}, map[string]interface{}{
		"username": username,
		"isAuth":   sess.IsAuth(),
	}, false)

	sess.Save(req, w)
	yRes := yoda.NewResponse(result.GetMeta(), result.GetData(), false)
	json.NewEncoder(w).Encode(yRes)
}
