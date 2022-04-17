package handlers

import (
	"log"
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/sockets"
	"github.com/gorilla/websocket"
)

// Establish handles creating a socket connection
func (env *Env) Establish(w http.ResponseWriter, req *http.Request) {
	store := env.createStore()
	defer store.Close()
	log.Printf("Store is %#v", store)

	sess, _ := env.sessionStore.Get(req, config.PostSessionName)
	if !sess.IsAuth() {
		yoda.SendClientError(w, yoda.ErrNotLoggedIn, http.StatusUnauthorized)
		sess.Save(req, w)
		return
	}

	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, req, nil)
	if err != nil {
		yoda.SendClientError(w, yoda.ErrCannotUpgrade, http.StatusBadRequest)
		sess.Save(req, w)
		return
	}

	id, _ := sess.Get("userid")
	client := sockets.NewClient(id.(string), conn, env.socketManager)

	go client.Read()
	go client.Write()
}
