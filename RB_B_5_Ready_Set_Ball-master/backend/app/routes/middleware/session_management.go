package middleware

import (
	"bytes"
	"crypto/sha1"
	"net/http"

	"log"

	"fmt"
	"time"

	"html"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/session"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
)

// func sessionMiddleware1(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		sess, err := env.sessionStore.Get(req, config.SessionName)
// 		if err != nil {
// 			yoda.SendClientError(w, yoda.ErrInvalidSession, http.StatusBadRequest)
// 			fmt.Printf("fake session %#v\n", err) //security log TODO: add more info like IP, location. user-agent, ...
// 			return
// 		}

// 		// Stop here if its Preflighted OPTIONS request
// 		if req.Method == "OPTIONS" {
// 			return
// 		}
// 		next.ServeHTTP(w, req)
// 	})
// }

// sessionMiddleware takes in a [session.Store]
// and returns a [http.Handler] function
// which takes in a http Handler style function
func sessionMiddleware(store session.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			sess, err := store.Get(req, config.PostSessionName) //Does this return an error if the user doesn't send it? That'll mean users can't sign in, fix it
			if err != nil {
				log.Printf("Security: Fatal: Error: Invalid session: %s\n%s\n", err.Error(), logAllDetails(req, sess)) //TODO: add more info location
				//destroy auth token and send new pre one
				sess.Destroy()
				sess.Save(req, w)
				sess, _ = store.New(req, config.SessionName)
				sess.Save(req, w)
				yoda.SendClientError(w, yerr.InvalidSession, http.StatusUnauthorized)
				return
			}

			// if the post is new, it means the user wasn't logged in and I prolly just genrated an auth sess so i'll change it
			if sess.IsNew() {
				sess, err = store.Get(req, config.SessionName)
				if err != nil {
					// security log and move on with new token
					log.Printf("Security: Error: Invalid session: %s\n%s\n", err.Error(), logAllDetails(req, sess)) //TODO: add more info location
				}
			}

			if isBlocked(sess) { //TODO: not tested manually
				log.Printf("Security: Fatal: Blocked session\n%s\n", logAllDetails(req, sess))
				sess.Destroy()
				sess.Save(req, w)
				yoda.SendClientError(w, yerr.InvalidSession, http.StatusBadRequest)
				return
			}

			//anomaly checks
			//TODO: consider making each of the anomaly chekcs a go routine and see how much faster it takes to complete
			// TODO: if you do that, use channels/select to stop all routines if one of them fails
			//user agent
			ua, ok := sess.Get("UserAgent")
			if !ok {
				rua := req.UserAgent()
				sess.Set("UserAgent", html.EscapeString(rua))
				ua, _ = sess.Get("UserAgent")
			}

			if agent, ok := ua.(string); !ok || agent != req.UserAgent() {
				log.Printf("Security: Fatal: User Agent change with the same session\n%s\n", logAllDetails(req, sess))
				sess.Destroy()
				sess.Save(req, w)
				yoda.SendClientError(w, yerr.InvalidSession, http.StatusBadRequest)
				return
			}

			if !config.Dev {
				//ip change
				ip, ok := sess.Get("IP")
				if !ok {
					addr := req.RemoteAddr
					sess.Set("IP", html.EscapeString(addr))
					ip, _ = sess.Get("IP")
				}
				if ipa, ok := ip.(string); !ok || ipa != req.RemoteAddr {
					log.Printf("Security: Fatal: IP address change with the same session\n%s\n", logAllDetails(req, sess))
					sess.Destroy()
					sess.Save(req, w)
					yoda.SendClientError(w, yerr.InvalidSession, http.StatusBadRequest)
					return
				}
			}

			//TODO check if expired, if gorilla doesn't already check

			next.ServeHTTP(w, req)
		})
	}
}

// logAllDetails returns a string containing logging information about the request and session
// which could then be passed into a logger
func logAllDetails(req *http.Request, sess session.Session) string {
	markIP(req, sess) //was a goroutine but whenever the sess got canclled, mark would panic because it marked too late
	return fmt.Sprintf("Time: %s\nUser-Agent: %s\nIP Address: %s\nMethod: %s\nPath: %s\nReferer: %s\nHost: %s\nNew Session-ID-hash: %x\nCookies: %s\n", time.Now().String(),
		req.UserAgent(),
		req.RemoteAddr,
		req.Method,
		req.RequestURI,
		req.Referer(),
		req.Host,
		sha1.Sum([]byte(sess.ID())),
		logCookies(req.Cookies()),
	)
}

// markRequest adds an IP to the database as bad so future requests from that IP will be blocked
func markIP(req *http.Request, sess session.Session) {
	logged, _ := sess.Get("logged")
	numLogged, _ := logged.(int)
	sess.Set("logged", numLogged+1)
	// Add IP to db of culprits
	// TODO: implement
}

// noteRequest just updates the session request number and
// if it's a login, signup, ... it also updates that number
func noteRequest(req *http.Request, sess session.Session) {
	// TODO: implement
	// should be done by each separte handler
}

// isBlocked returns true if a session is blocked from too many failed attempts
func isBlocked(sess session.Session) bool {
	logged, _ := sess.Get("logged")
	numLogged, _ := logged.(int)
	return numLogged >= config.MaxSessionAttempts
	// TODO: this should be done based on the ip from a culprit db
}

// logCookies returns a readable array of cookies in a Name - Value pair
func logCookies(cookies []*http.Cookie) string {
	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for _, cookie := range cookies {
		buffer.WriteString("\tName: ")
		buffer.WriteString(cookie.Name)
		buffer.WriteString(", Value hash: ")
		hash := sha1.Sum([]byte(cookie.Value))
		buffer.WriteString(fmt.Sprintf("%x", hash))
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}

//TODOS
// Block IPs that send two consecutive expired/invalid time stamps
