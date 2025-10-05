package service

import (
	e "calendar/internal/entity"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type CalendarService struct {
	EventInfo map[uuid.UUID][]e.Event
}

func NewCalendarService() *CalendarService {
	return &CalendarService{
		EventInfo: make(map[uuid.UUID][]e.Event),
	}
}

func (c *CalendarService) AddNewEvent(ctx context.Context, event e.NewEvent) (uuid.UUID, error) {
	_, err := time.Parse("2006-January-2", dateToString(event))
	if err != nil {
		log.Printf("invalid date %d-%s-%d: %v\n use format 2006-January-2", event.Date.Year, event.Date.Month.String(), event.Date.Day, err)
		return uuid.Nil, err
	}

	eventID := uuid.New()

	c.EventInfo[event.UserID] = append(c.EventInfo[event.UserID], e.Event{
		EventID:    eventID,
		Date:       event.Date,
		EventInfo:  event.EventInfo,
		CreateDate: time.Now(),
	})

	return eventID, nil
}

func (c *CalendarService) UpdateEvent(ctx context.Context, newEvent e.NewEvent, eventID string) error {
	_, err := time.Parse("2006-January-2", dateToString(newEvent))
	if err != nil {
		log.Printf("invalid date %d-%s-%d: %v\n use format 2006-January-2", newEvent.Date.Year, newEvent.Date.Month.String(), newEvent.Date.Day, err)
		return err
	}

	id, err := uuid.Parse(eventID)
	if err != nil {
		log.Printf("invalid event id: %s", eventID)
		return err
	}

	for i := range c.EventInfo[newEvent.UserID] {
		if id == c.EventInfo[newEvent.UserID][i].EventID {
			c.EventInfo[newEvent.UserID][i].Date = newEvent.Date
			c.EventInfo[newEvent.UserID][i].EventInfo = newEvent.EventInfo
			break
		}
	}

	return nil
}

func (c *CalendarService) DeleteEvent(ctx context.Context, userIDstr string, eventIDstr string) error {
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		log.Printf("invalid user id: %v", err)
		return fmt.Errorf("invalid user id: %v", err)
	}

	eventID, err := uuid.Parse(eventIDstr)
	if err != nil {
		log.Printf("invalid event id: %v", err)
		return fmt.Errorf("invalid event id: %v", err)
	}

	for i := range c.EventInfo[userID] {
		if c.EventInfo[userID][i].EventID == eventID {
			c.EventInfo[userID][i] = e.Event{}
			return nil
		}
	}
	return fmt.Errorf("non-existent ID")
}

func (c *CalendarService) GetEventsForDay(stringID string, date string) ([]e.Event, bool) {
	t, err := time.Parse("2006-January-2", date)
	if err != nil {
		log.Printf("invalid date %s: %v\n use format 2006-January-2", date, err)
		return nil, false
	}

	id, err := uuid.Parse(stringID)
	if err != nil {
		log.Printf("invalid id: %v", err)
		return nil, false
	}
	var events []e.Event
	var status bool
	targetDate := e.Date{Year: t.Year(), Month: t.Month(), Day: t.Day()}

	for _, event := range c.EventInfo[id] {
		if event.Date == targetDate {
			events = append(events, event)
			status = true
		}
	}

	return events, status
}

func (c *CalendarService) GetEventsForWeek(userIDstr string, date string) ([]e.Event, bool) {
	t, err := time.Parse("2006-January-2", date)
	if err != nil {
		log.Printf("invalid date %s: %v\n use format 2006-January-2", date, err)
		return nil, false
	}
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		log.Printf("invalid id: %v", err)
		return nil, false
	}

	var events []e.Event
	status := false
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := t.AddDate(0, 0, -weekday+1)

	var week []e.Date
	for i := 0; i < 7; i++ {
		tmp := monday.AddDate(0, 0, i)
		targetDate := e.Date{Year: tmp.Year(), Month: tmp.Month(), Day: tmp.Day()}
		week = append(week, targetDate)
	}

	for i := range c.EventInfo[userID] {
		for j := range week {
			ev := c.EventInfo[userID][i].Date
			wk := week[j]
			if ev.Year == wk.Year && ev.Month == wk.Month && ev.Day == wk.Day {
				events = append(events, c.EventInfo[userID][i])
				status = true
			}
		}
	}

	return events, status
}

func (c *CalendarService) GetEventsForMonth(userIDstr string, date string) ([]e.Event, bool) {
	t, err := time.Parse("2006-January-2", date)
	if err != nil {
		log.Printf("invalid date %s: %v\n use format 2006-January-2", date, err)
		return nil, false
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		log.Printf("invalid user id: %v", err)
		return nil, false
	}
	var events []e.Event
	var status bool
	targetDate := e.Date{Year: t.Year(), Month: t.Month()}

	for _, event := range c.EventInfo[userID] {
		if event.Date.Month == targetDate.Month && event.Date.Year == targetDate.Year {
			events = append(events, event)
			status = true
		}
	}

	return events, status
}

func dateToString(event e.NewEvent) string {
	return strconv.Itoa(event.Date.Year) + "-" + event.Date.Month.String() + "-" + strconv.Itoa(event.Date.Day)
}
