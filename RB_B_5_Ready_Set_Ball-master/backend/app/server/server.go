package server

import (
	"log"
	"net/http"
	"strconv"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/database"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/routes"
)

/**TODO:
- Include logger
- Enable HTTPS
*/

// Server represents the server which listens for connections when started
type Server struct {
	Hostname  string `json:"hostname"`  // Server name
	UseHTTP   bool   `json:"UseHTTP"`   // Listen on HTTP
	UseHTTPS  bool   `json:"UseHTTPS"`  // Listen on HTTPS
	HTTPPort  int    `json:"HTTPPort"`  // HTTP port
	HTTPSPort int    `json:"HTTPSPort"` // HTTPS port
	CertFile  string `json:"CertFile"`  // HTTPS certificate
	KeyFile   string `json:"KeyFile"`   // HTTPS private key
	Handler   http.Handler
}

// New creates a new server using the [config] map
// func (s *Server) New(config map[string]string) {
// 	// config["o"] = "p"
// 	s.Handler = routes.Load()
// 	// s.HTTPPort = 4444
// 	// s.UseHTTP = true
// }

// Start fires a listener and starts the server on the specified port
// using HTTPS if [Server.UseHTTPS] is true else it uses HTTP
// It creates the database connection
func (s *Server) Start() {
	var db *database.Database
	var err error

	if config.Beta {
		db, err = database.ConnectWithInfo(config.MongoDBHosts)
	} else {
		db, err = database.Connect("localhost")
	}
	if err != nil {
		log.Fatal(err) //TODO: panic and recover
	}
	defer db.Close()

	// attaches the server handler
	if s.Handler == nil {
		s.Handler = routes.Load(db.NewStoreFactory())
	}

	if s.UseHTTPS {
		s.startHTTPS()
	} else {
		s.startHTTP()
	}
}

func (s *Server) startHTTP() {
	//TODO: Use hostanme too
	log.Printf("Server started on %d\n", s.HTTPPort)
	log.Fatal(http.ListenAndServe(s.address(), s.Handler))
}

func (s *Server) startHTTPS() {
	//TODO: Complete
	log.Fatal(http.ListenAndServeTLS(s.address(), s.CertFile, s.KeyFile, s.Handler))
}

func (s *Server) address() string {
	if s.UseHTTPS {
		return ":" + strconv.Itoa(s.HTTPSPort)
	}
	return ":" + strconv.Itoa(s.HTTPPort)
}
