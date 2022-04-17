package routes

import (
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/socket_handler"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/sockets"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/session"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/routes/handlers"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/routes/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func getRouter(env *handlers.Env) *httprouter.Router {
	router := httprouter.New()
	handler := middleware.Handler
	//TODO: I don't think I need alice on any of these. Infact I'll only need alice if I was adding specific middleware
	router.GET("/people", handler(alice.New().ThenFunc(env.GetUsers)))

	// HTTP Errors
	router.NotFound = alice.New().ThenFunc(env.NotFound)
	router.MethodNotAllowed = alice.New().ThenFunc(env.MethodNotAllowed)

	// Game
	router.GET("/games/l/lng/:lng/lat/:lat", handler(alice.New().ThenFunc(env.GetGamesNearLocation)))
	router.GET("/game", handler(alice.New().ThenFunc(env.GetGame)))
	router.GET("/game/g/:id", handler(alice.New().ThenFunc(env.GetGame)))
	router.POST("/game/exit", handler(alice.New().ThenFunc(env.ExitGame)))
	router.POST("/game/join/:mode", handler(alice.New().ThenFunc(env.JoinGame)))
	router.POST("/game/rate", handler(alice.New().ThenFunc(env.RateGame)))

	// User
	router.POST("/user/:command/:argument", handler(alice.New().ThenFunc(env.GetUser)))
	router.POST("/upload", handler(alice.New().ThenFunc(env.UploadAvatar)))

	// Remove entity
	router.POST("/remove/:item", handler(alice.New().ThenFunc(env.Remove)))

	// Create entity
	router.POST("/create/:entity", handler(alice.New().ThenFunc(env.CreateEntity)))

	// Edit entity
	router.POST("/edit/:entity", handler(alice.New().ThenFunc(env.Edit)))

	// Invites
	router.POST("/invite/m/:mode/t/:type", handler(alice.New().ThenFunc(env.Invites)))

	// Authentication
	router.POST("/login", handler(alice.New().ThenFunc(env.Login)))
	router.GET("/logout", handler(alice.New().ThenFunc(env.Logout)))

	// Status
	router.GET("/status", handler(alice.New().ThenFunc(env.Status)))

	// Socket
	router.GET("/establish", handler(alice.New().ThenFunc(env.Establish)))

	/*
		Imagine a case where Env was this:
		env.GetNotes returns a type Handler
		e.g.
		type Handler struct {
			httpmethods string //e.g "POST", "GET"
			path string e.g "/people/:id"
			handler httphandler e.g a normal golang handler
		}
		Then we could do
		router.Handle(env.GetUsers.httpmethods, env.GetUsers.path, env.GetUsers.handler)

		We could even take it further and put all routes in an error then a for loop to go through all and set up the routes
	*/

	return router
}

// Load returns the routes and middleware
// This also handles the inclusion of a database session and browser session
func Load(storeFactory model.StoreFactory) http.Handler { //TODO: complete this and rename New-ish?

	// MongoDB store for storing the session
	mongoStore := storeFactory() //NOTE: possible leak here from not closing the store, but it stores our sessions, do we really want to close it?
	sessStore := session.NewMongoSessionStore(mongoStore.GetCollection(config.SessionStoreName), []byte(config.SessionAuth))

	// setup socket listener
	socketManager := sockets.NewManager(sockethandler.NewSocketHandler())
	go socketManager.Start()

	// Setup environment with the datastore and session store
	env := handlers.NewEnv(storeFactory, sessStore, socketManager)

	return middleware.ApplyMiddleware(getRouter(env), sessStore)
}
