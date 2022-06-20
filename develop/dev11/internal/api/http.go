package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	logger "httpserv/internal/logger"
)

func requestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startedTime := time.Now()
		logger.Log(r, "started")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)

		endTime := time.Now()
		logText := fmt.Sprintf("completed for %v", endTime.Sub(startedTime))
		logger.Log(r, logText)
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rlog := logger.RequestLog{
			Time:       time.Now(),
			RemoteAddr: r.RemoteAddr,
			Method:     r.Method,
			Header:     r.UserAgent(),
		}
		ctx := context.WithValue(context.Background(), "request", rlog)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func StartHTTP() error {
	http.HandleFunc("/create_event", loggerMiddleware(requestMiddleware(CreateEvent)))
	http.HandleFunc("/update_event", loggerMiddleware(requestMiddleware(UpdateEvent)))
	http.HandleFunc("/delete_event", loggerMiddleware(requestMiddleware(DeleteEvent)))
	http.HandleFunc("/events_for_day", loggerMiddleware(requestMiddleware(GetEventsForDay)))
	http.HandleFunc("/events_for_week", loggerMiddleware(requestMiddleware(GetEventsForWeek)))
	http.HandleFunc("/events_for_month", loggerMiddleware(requestMiddleware(GetEventsForMonth)))
	http.ListenAndServe("127.0.0.1:8080", nil)
	return nil
}
