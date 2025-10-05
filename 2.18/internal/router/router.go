package router

import (
	h "calendar/internal/handler"
	"calendar/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(h *h.CalendarHandler) *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.Logger)

	r.HandleFunc("/create_event", h.CreateCalendar).Methods("POST")
	r.HandleFunc("/update_event/{eventID}", h.UpdateEvent).Methods("POST")
	r.HandleFunc("/delete_event/{eventID}", h.DeleteEvent).Methods("POST")
	r.HandleFunc("/events_for_day", h.GetEventsForDay).Methods("GET")
	r.HandleFunc("/events_for_week", h.GetEventsForWeek).Methods("GET")
	r.HandleFunc("/events_for_month", h.GetEventsForMonth).Methods("GET")

	return r
}
