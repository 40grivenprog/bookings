package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/40grivenprog/bookings/internal/config"
	"github.com/40grivenprog/bookings/internal/driver"
	"github.com/40grivenprog/bookings/internal/forms"
	"github.com/40grivenprog/bookings/internal/helpers"
	"github.com/40grivenprog/bookings/internal/models"
	"github.com/40grivenprog/bookings/internal/render"
	"github.com/40grivenprog/bookings/internal/repository"
	"github.com/40grivenprog/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DataBaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "home.page.tmpl", &models.TemplateData{}, r)
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send data to the template
	render.Template(w, "about.page.tmpl", &models.TemplateData{}, r)
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data["reservation"] = res

	render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	}, r)
}

// PostReservation handles posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	start_date, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	end_date, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: start_date,
		EndDate:   end_date,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	form.MinLength("first_name", 3, r)
	form.Required("first_name", "last_name", "email", "phone")
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	}
	newReservationId, err := m.DB.InsertReservation(reservation)

	if err != nil {
		helpers.ServerError(w, err)
	}

	restriction := models.RoomRestriction{
		StartDate:     start_date,
		EndDate:       end_date,
		RoomID:        roomID,
		ReservationID: newReservationId,
		RestrictionID: 1,
	}
	err = m.DB.InsertRoomRestriction(restriction)

	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "generals.page.tmpl", &models.TemplateData{}, r)
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "majors.page.tmpl", &models.TemplateData{}, r)
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "search-availability.page.tmpl", &models.TemplateData{}, r)
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")

	end := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, "choose-room.page.tmpl", &models.TemplateData{Data: data}, r)
}

type jsonResponse struct {
	OK        bool   `json:"ok"` // need for marshall. it checks this values for json fields
	Message   string `json:"message"`
	RoomId    string `json:"room_id"`
	EndDate   string `json:"end_date"`
	StartDate string `json:"start_date"`
}

// AvailabilityJSON renders the search availability page
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomId, _ := strconv.Atoi(r.Form.Get("room_id"))
	log.Println(roomId)

	available, _ := m.DB.SearchAvailabilityByDateByRoomId(startDate, endDate, roomId)
	log.Println(available)

	resp := jsonResponse{
		OK:        available,
		Message:   "OK!",
		RoomId:    r.Form.Get("room_id"),
		StartDate: sd,
		EndDate:   ed,
	}

	out, err := json.MarshalIndent(resp, "", "    ")

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	log.Println(string(out))

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "contact.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		log.Println("Error!")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, "reservation-summary.page.tmpl", &models.TemplateData{Data: data}, r)
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomId
	room, err := m.DB.GetRoomByID(roomId)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room = room

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}


func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	layout := "2006-01-02"
	startDate := r.URL.Query().Get("s")
	endDate := r.URL.Query().Get("e")

	start_date, err := time.Parse(layout, startDate)
	if err != nil {
		helpers.ServerError(w, err)
	}
	end_date, err := time.Parse(layout, endDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	var res models.Reservation
	
	res.RoomID = roomID
	res.StartDate = start_date
	res.EndDate = end_date

	room, err := m.DB.GetRoomByID(roomID)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
