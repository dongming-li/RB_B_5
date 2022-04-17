package middleware

import (
	"net/http"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/session"
	"github.com/justinas/alice"
)

/**TODO: Implement the following middleware
1. CORS
2. CSRF
3. Oauth
4. Authentication
5. URI handling
6. Logging
*/

// Middleware is a type to represent the middleware handler function
// type Middleware func(http.Handler) http.Handler

// ApplyMiddleware is the overall middleware for all requests
func ApplyMiddleware(h http.Handler, store session.Store) http.Handler {
	// Add more middleware here
	return alice.New(loggingHandler, sessionMiddleware(store), corsMiddleware).Then(h)
}
