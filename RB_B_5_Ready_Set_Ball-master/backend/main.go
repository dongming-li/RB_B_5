package main

import "git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/server"

func main() {
	// Use all CPU cores
	// runtime.GOMAXPROCS(runtime.NumCPU())

	s := &server.Server{
		UseHTTP:  true,
		HTTPPort: 4444,
	}

	s.Start()
}
