package tests

import (
	e "calendar/internal/entity"
	s "calendar/internal/service"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestСreateEvent(t *testing.T) {
	uuidUser, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		log.Printf("error parse uuid: %v\n", err)
		return
	}

	date := e.Date{
		Year:  2025,
		Month: time.October,
		Day:   12,
	}
	event := e.NewEvent{
		UserID:    uuidUser,
		Date:      date,
		EventInfo: "New event",
	}
	svc := s.NewCalendarService()

	_, err = svc.AddNewEvent(context.Background(), event)
	if err != nil {
		fmt.Println(err)
	}

	events, ok := svc.GetEventsForDay(uuidUser.String(), "2025-October-12")
	if !ok {
		log.Printf("No events")
		return
	}

	for _, e := range events {
		log.Println(e)
	}
}

func TestUpdateEvent(t *testing.T) {
	uuidUser, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		log.Printf("error parse uuid: %v\n", err)
		return
	}

	date := e.Date{
		Year:  2025,
		Month: time.October,
		Day:   12,
	}
	event := e.NewEvent{
		UserID:    uuidUser,
		Date:      date,
		EventInfo: "New event",
	}
	svc := s.NewCalendarService()
	ctx := context.Background()
	eventID, err := svc.AddNewEvent(ctx, event)
	if err != nil {
		fmt.Println(err)
	}

	events, ok := svc.GetEventsForDay(uuidUser.String(), "2025-October-12")
	if !ok {
		log.Printf("No events")
		return
	}

	for _, e := range events {
		log.Println(e)
	}

	newDate := e.Date{
		Year:  2025,
		Month: time.October,
		Day:   15,
	}
	newEvent := e.NewEvent{
		UserID:    uuidUser,
		Date:      newDate,
		EventInfo: "Updated event",
	}

	err = svc.UpdateEvent(ctx, newEvent, eventID.String())
	if err != nil {
		log.Print(err)
		return
	}

	events, ok = svc.GetEventsForDay(uuidUser.String(), "2025-October-15")
	if !ok {
		log.Printf("No events")
		return
	}

	for _, e := range events {
		log.Println(e)
	}
}

func TestDeleteEvent(t *testing.T) {
	uuidUser, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		log.Printf("error parse uuid: %v\n", err)
		return
	}

	date := e.Date{
		Year:  2025,
		Month: time.October,
		Day:   12,
	}
	event := e.NewEvent{
		UserID:    uuidUser,
		Date:      date,
		EventInfo: "New event",
	}
	svc := s.NewCalendarService()
	ctx := context.Background()
	eventID, err := svc.AddNewEvent(ctx, event)
	if err != nil {
		fmt.Println(err)
	}

	events, ok := svc.GetEventsForDay(uuidUser.String(), "2025-October-12")
	if !ok {
		log.Printf("No events")
		return
	}

	for _, e := range events {
		log.Println(e)
	}

	err = svc.DeleteEvent(ctx, uuidUser.String(), eventID.String())
	if err != nil {
		log.Printf("Error delete event: %v", err)
		return
	}

	events, ok = svc.GetEventsForDay(uuidUser.String(), "2025-October-12")
	if !ok {
		log.Printf("No events")
		return
	}

	for _, e := range events {
		log.Println(e)
	}
}

func TestGet(t *testing.T) {
	uuidUser, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		log.Printf("error parse uuid: %v\n", err)
		return
	}

	date1 := e.Date{
		Year:  2025,
		Month: time.October,
		Day:   12,
	}
	event1 := e.NewEvent{
		UserID:    uuidUser,
		Date:      date1,
		EventInfo: "New event1",
	}
	svc := s.NewCalendarService()
	ctx := context.Background()
	_, err = svc.AddNewEvent(ctx, event1)
	if err != nil {
		fmt.Println(err)
	}

	date2 := e.Date{
		Year:  2025,
		Month: time.November,
		Day:   21,
	}
	event2 := e.NewEvent{
		UserID:    uuidUser,
		Date:      date2,
		EventInfo: "New event2",
	}

	_, err = svc.AddNewEvent(ctx, event2)
	if err != nil {
		fmt.Println(err)
	}

	date3 := e.Date{
		Year:  2025,
		Month: time.October,
		Day:   10,
	}
	event3 := e.NewEvent{
		UserID:    uuidUser,
		Date:      date3,
		EventInfo: "New event3",
	}

	_, err = svc.AddNewEvent(ctx, event3)
	if err != nil {
		fmt.Println(err)
	}

	date4 := e.Date{
		Year:  2025,
		Month: time.October,
		Day:   19,
	}
	event4 := e.NewEvent{
		UserID:    uuidUser,
		Date:      date4,
		EventInfo: "New event3",
	}

	_, err = svc.AddNewEvent(ctx, event4)
	if err != nil {
		fmt.Println(err)
	}

	events, ok := svc.GetEventsForWeek(uuidUser.String(), "2025-October-12") // get week events
	if !ok {
		log.Printf("No events")
		return
	}

	log.Println("Result get week events:")
	for _, e := range events {
		log.Println(e)
	}

	events, ok = svc.GetEventsForMonth(uuidUser.String(), "2025-October-12") // get month events
	if !ok {
		log.Printf("No events")
		return
	}
	log.Println("Result get month events:")
	for _, e := range events {
		log.Println(e)
	}

}
