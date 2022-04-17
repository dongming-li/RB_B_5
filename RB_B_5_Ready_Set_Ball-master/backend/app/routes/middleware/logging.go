package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func loggingHandler(h http.Handler) http.Handler {
	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to use open log buffer:", err)
		return h
	}
	return handlers.LoggingHandler(logFile, h)
}
