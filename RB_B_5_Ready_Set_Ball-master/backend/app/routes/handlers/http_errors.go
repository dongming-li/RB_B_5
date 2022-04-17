package handlers

import (
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// NotFound is the route for a 404 Not found Error
// It returns {meta: {error: true}, result: {error: NotFound}}
func (env *Env) NotFound(w http.ResponseWriter, req *http.Request) {
	sess, _ := env.sessionStore.Get(req, config.PostSessionName)
	sess.Save(req, w)
	yoda.SendClientError(w, yoda.ErrResourceNotFound, http.StatusNotFound)
}

// MethodNotAllowed is the route for a 405 Method Not Allowed Error
// It returns {meta: {error: true}, result: {error: NotAllowed}}
func (env *Env) MethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	sess, _ := env.sessionStore.Get(req, config.PostSessionName)
	sess.Save(req, w)
	yoda.SendClientError(w, yoda.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
}
