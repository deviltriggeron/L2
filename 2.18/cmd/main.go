package main

import (
	h "calendar/internal/handler"
	r "calendar/internal/router"
	s "calendar/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	svc := s.NewCalendarService()
	handler := h.NewCalendarHandler(svc)
	router := r.NewRouter(handler)

	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port := os.Getenv("LOCALHOST_PORT")
	srv := http.Server{
		Addr:    port,
		Handler: router,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("Listen and running %s:\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	srv.Shutdown(ctx)

	wg.Wait()
}
