package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go"
)

const defaultSleepTime = time.Second * 2

func main() {
	httpPort := flag.Int("http.port", 4000, "HTTP Port to run server on")
	mongoDSN := flag.String("mongo.dsn", "localhost:27017", "DSN for mongoDB server")

	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appID := os.Getenv("PUSHER_APP_ID")
	appKey := os.Getenv("PUSHER_APP_KEY")
	appSecret := os.Getenv("PUSHER_APP_SECRET")
	appCluster := os.Getenv("PUSHER_APP_CLUSTER")
	appIsSecure := os.Getenv("PUSHER_APP_SECURE")

	var isSecure bool
	if appIsSecure == "1" {
		isSecure = true
	}

	client := &pusher.Client{
		AppId:   appID,
		Key:     appKey,
		Secret:  appSecret,
		Cluster: appCluster,
		Secure:  isSecure,
		HttpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	mux := chi.NewRouter()

	log.Println("Connecting to MongoDB")
	m, err := newMongo(*mongoDSN)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to MongoDB")

	mux.Use(analyticsMiddleware(m, client))

	var once sync.Once
	var t *template.Template

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	fileServer(mux, "/static", http.Dir(filesDir))

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {

		once.Do(func() {
			tem, err := template.ParseFiles("index.html")
			if err != nil {
				log.Fatal(err)
			}

			t = tem.Lookup("index.html")
		})

		t.Execute(w, nil)
	})

	mux.Get("/api/analytics", analyticsAPI(m))
	mux.Get("/wait/{seconds}", waitHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), mux))
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}

	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func analyticsAPI(m mongo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		errFn := func(w http.ResponseWriter) {

			json.NewEncoder(w).Encode(&struct {
				Message   string `json:"message"`
				TimeStamp int64  `json:"timestamp"`
			}{
				Message:   "An error occurred while fetching analytics data",
				TimeStamp: time.Now().Unix(),
			})
		}

		totalRequests, err := m.Count()
		if err != nil {
			errFn(w)
			return
		}

		stats, err := m.StatsPerRoute()
		if err != nil {
			errFn(w)
			return
		}

		reqsPerDay, err := m.RequestsPerDay()
		if err != nil {
			errFn(w)
			return
		}

		reqsPerHour, err := m.RequestsPerHour()
		if err != nil {
			errFn(w)
			return
		}

		avgResponseTime, err := m.AverageResponseTime()
		if err != nil {
			errFn(w)
			return
		}

		type Data struct {
			AverageResponseTime float64          `json:"average_response_time"`
			StatsPerRoute       []statsPerRoute  `json:"stats_per_route"`
			RequestsPerDay      []requestsPerDay `json:"requests_per_day"`
			RequestsPerHour     []requestsPerDay `json:"requests_per_hour"`
			TotalRequests       int              `json:"total_requests"`
		}

		json.NewEncoder(w).Encode(Data{
			AverageResponseTime: avgResponseTime,
			StatsPerRoute:       stats,
			RequestsPerDay:      reqsPerDay,
			RequestsPerHour:     reqsPerHour,
			TotalRequests:       totalRequests,
		})
	}
}

func analyticsMiddleware(m mongo, client *pusher.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			startTime := time.Now()

			defer func() {
				switch r.URL.String() {
				case "/analytics", "/api/analytics", "/favicon.ico", "/serviceworker.js":

				default:
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

					client.Trigger("analytics-dashboard", "data", data)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
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
