package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// Logout handles request to destroy an auth
func (env *Env) Logout(w http.ResponseWriter, req *http.Request) { // /login Potentially could login/:source where source is an integer to represent where the login is coming from

	sess, _ := env.sessionStore.Get(req, config.PostSessionName)
	// if there was an error, the req would've been rejected by sess middlware
	//and also we're logging out, we don't care
	//This if/else block is for logging purposes only
	if sess.IsAuth() {
		username, _ := sess.Get("username")
		log.Println("Logging out", username)
	} else {
		log.Println("Destroying sessions with no user to logout")
	}
	sess.Destroy()
	sess.Save(req, w)

	sess, _ = env.sessionStore.Get(req, config.SessionName)
	//we're logging out, we don't care if there was an error. TODO: should we?
	sess.Destroy()
	sess.Save(req, w)
	//TODO: operation.Record("logout")

	yRes := yoda.NewResponse(map[string]string{}, map[string]string{"logout": "successful"}, false) //TODO use one argument newresponse
	json.NewEncoder(w).Encode(yRes)
}
