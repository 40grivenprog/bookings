package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/40grivenprog/bookings/internal/config"
	"github.com/40grivenprog/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)



var session *scs.SessionManager
var testApp config.AppConfig


func TestMain(m *testing.M) {
		// what i am going to put in the session
		gob.Register(models.Reservation{})
		// change this to true when in production
		testApp.InProduction = false
	
		// set up the session
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = testApp.InProduction
	
		testApp.Session = session

		app = &testApp

	os.Exit(m.Run())
}


type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
