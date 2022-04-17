package handlers

import (
	"context"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/sockets"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/session"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/routes/middleware"
)

// Env is a wrapper for the genral request handling and contains a store instance factory
type Env struct {
	createStore   model.StoreFactory
	sessionStore  session.Store
	socketManager *sockets.Manager
}

// NewEnv returns a new [Env] with the specified [model.StoreFactory]
func NewEnv(storeFactory model.StoreFactory, sessStore session.Store, socketManager *sockets.Manager) *Env {
	return &Env{
		createStore:   storeFactory,
		sessionStore:  sessStore,
		socketManager: socketManager,
	}
}

func retrieveParamsFromContext(ctx context.Context) map[string]string {
	//TODO: check for error
	var key middleware.ParamKey = "params"
	return ctx.Value(key).(map[string]string)
}
