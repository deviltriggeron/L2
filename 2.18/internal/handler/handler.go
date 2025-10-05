package handler

import (
	e "calendar/internal/entity"
	s "calendar/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CalendarHandler struct {
	svc *s.CalendarService
}

func NewCalendarHandler(svc *s.CalendarService) *CalendarHandler {
	return &CalendarHandler{
		svc: svc,
	}
}

func (h *CalendarHandler) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	var newEvent e.NewEvent
	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		respondJSON(w, http.StatusBadRequest, "")
		return
	}

	eventID, err := h.svc.AddNewEvent(r.Context(), newEvent)
	if err != nil {
		respondJSON(w, http.StatusServiceUnavailable, "")
		return
	}
	respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
		"UserID":                 newEvent.UserID,
		"Success create eventID": eventID,
	})
}

func (h *CalendarHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent e.NewEvent
	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		respondJSON(w, http.StatusBadRequest, "")
		return
	}

	vars := mux.Vars(r)
	id := vars["eventID"]

	if err := h.svc.UpdateEvent(r.Context(), newEvent, id); err != nil {
		respondJSON(w, http.StatusServiceUnavailable, "")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"Success update event id": id,
	})
}

func (h *CalendarHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var deletedUserID string
	if err := json.NewDecoder(r.Body).Decode(&deletedUserID); err != nil {
		respondJSON(w, http.StatusBadRequest, "")
		return
	}

	vars := mux.Vars(r)
	id := vars["eventID"]
	err := h.svc.DeleteEvent(r.Context(), deletedUserID, id)
	if err != nil {
		respondJSON(w, http.StatusServiceUnavailable, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"Event succesfully deleted": id})
}

func (h *CalendarHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()

	id := params.Get("id")
	date := params.Get("date")

	events, ok := h.svc.GetEventsForDay(id, date)
	if ok {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"Events: ": events,
		})
	} else {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"Not event in: ": date,
		})
	}
}

func (h *CalendarHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()

	id := params.Get("id")
	date := params.Get("date")

	events, ok := h.svc.GetEventsForWeek(id, date)
	if ok {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"Events: ": events,
		})
	} else {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"Not event in: ": date,
		})
	}
}

func (h *CalendarHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	id := params.Get("id")
	date := params.Get("date")

	events, ok := h.svc.GetEventsForMonth(id, date)
	if ok {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"Events: ": events,
		})
	} else {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"Not event in: ": date,
		})
	}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
