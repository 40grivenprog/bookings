package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)


// adds CSRF protection for all POST
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// Loads session and saves it on every request(as cookie)
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
