package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// Login handles request to authenticate a user
func (env *Env) Login(w http.ResponseWriter, req *http.Request) { // /login Potentially could login/:source where source is an integer to represent where the login is coming from
	store := env.createStore()
	defer store.Close()
	log.Printf("Store is %#v", store)

	//session
	sess, err := env.sessionStore.Get(req, config.PostSessionName)
	// if there was an error, the req would've been rejected by sess middlware
	if sess.IsNew() {
		sess, err = env.sessionStore.Get(req, config.SessionName)
		if err != nil {
			yoda.SendClientError(w, yoda.ErrInvalidSession, http.StatusBadRequest)
			fmt.Printf("fake session %#v\n", err) //seurity log TODO: add more info like IP, location. user-agent, ...
			return
		}
	} else {
		//TODO: This error is sent even when a wrong username/password is sent but a correct session is given
		//		Consider if sess.IsNew() && sess.Get("username") == req.Result["username"] {} else {code below}
		// Check if user is already logged in
		if sess.IsAuth() {
			yoda.SendClientError(w, yoda.ErrNoLogin, http.StatusBadRequest)
			return
		}
		// else {
		//should be impossible to hit this else
		//TODO: destroy auth token, generate preauth token and send it
		// sess.Destroy()
		// sess.Save(req, w)
		// }
	}

	yReq, err := yoda.NewRequest(req, nil)
	if err != nil {
		yoda.SendClientError(w, err, 0) //TODO: add code
		return
	}

	result, err := operations.Do("login", store, *yReq) //authenticate vs login?
	if err != nil {
		sess.Save(req, w)
		yoda.SendClientError(w, err, 0) //TODO: add code
		return
	}

	sess.Destroy()
	sess.Save(req, w) //TODO: check err?
	sess, err = env.sessionStore.New(req, config.PostSessionName)
	if err != nil {
		//Note: if i save here, it will  be an unauthorized auth token that isn't new, that's just weird
		//should prolly send a preauth token
		yoda.SendClientError(w, err, 500) //TODO: check this right
		return
	}

	//save session
	err = sess.Assign(result) //What's happening here?
	if err != nil {
		//Note: if i save here, it will  be an unauthorized auth token that isn't new, that's just weird
		//should prolly send a preauth token
		yoda.SendClientError(w, err, 500) //TODO: add error code
		return
	}

	sess.Save(req, w)                                                   //TODO: check err?
	yRes := yoda.NewResponse(result.GetMeta(), result.GetData(), false) //TODO use one argument newresponse
	json.NewEncoder(w).Encode(yRes)
}
