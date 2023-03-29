package main

import (
	"fmt"
	"log"
	"github.com/40grivenprog/bookings/pkg/config"
	"github.com/40grivenprog/bookings/pkg/handlers"
	"github.com/40grivenprog/bookings/pkg/render"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber string = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // Persist when browser closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session



	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache!")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %v", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	_ = srv.ListenAndServe()
}
