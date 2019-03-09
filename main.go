package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

const defaultSleepTime = time.Second * 2

func main() {
	httpPort := flag.Int("http.port", 4000, "HTTP Port to run server on")
	mongoDSN := flag.String("mongo.dsn", "localhost:27017", "DSN for mongoDB server")

	flag.Parse()

	mux := chi.NewRouter()

	log.Println("Connecting to MongoDB")
	m, err := newMongo(*mongoDSN)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to MongoDB")

	mux.Use(analyticsMiddleware(m))
	mux.Get("/analytics", displayAnalytics)
	mux.Get("/wait/{seconds}", waitHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), mux))
}

func analyticsMiddleware(m mongo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			startTime := time.Now()

			defer func() {
				if r.URL.String() != "/analytics" {
					data := requestAnalytics{
						URL:         r.URL.String(),
						Method:      r.Method,
						RequestTime: time.Now().Unix() - startTime.Unix(),
						Day:         startTime.Weekday(),
						Hour:        startTime.Hour(),
					}

					if err := m.Write(data); err != nil {
						log.Println(err)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func displayAnalytics(w http.ResponseWriter, r *http.Request) {

}

func waitHandler(w http.ResponseWriter, r *http.Request) {
	var sleepTime = defaultSleepTime

	secondsToSleep := chi.URLParam(r, "seconds")
	n, err := strconv.Atoi(secondsToSleep)
	if err == nil && n >= 2 {
		sleepTime = time.Duration(n) * time.Second
	} else {
		n = 2
	}

	log.Printf("Sleeping for %d seconds", n)
	time.Sleep(sleepTime)
	w.Write([]byte(`Done`))
}
