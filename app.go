package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/alvintzz/alert-thread/internal/handler"
	"github.com/alvintzz/alert-thread/internal/repository/notification/slack"
	"github.com/alvintzz/alert-thread/internal/repository/storage/gmap"
	"github.com/alvintzz/alert-thread/internal/usecase"
)

// Config is main configuraton for slack-alert service
type Config struct {
	Server Server `json:"server"`
	Log    Log    `json:"log"`
	Slack  Slack  `json:"slack"`
}

// Server defines server config for http server
type Server struct {
	Port         string        `json:"port"`
	WriteTimeout time.Duration `json:"write_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
}

// Log defines log configuration of the service
type Log struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// Slack defines slack configuration
type Slack struct {
	Token string `json:"token"`
}

func main() {
	// Get Flag parameter from user
	var configFile string
	flag.StringVar(&configFile, "config.file", "config/config.json", "Path of config location")

	// Parse input flag
	flag.Parse()

	// Initialize Config
	config, err := readConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Log Level Logging
	err = initLog(config.Log.Level, config.Log.Format)
	if err != nil {
		log.Fatal("Failed to initialize Logger", err)
	}

	incidentStorage, err := gmap.NewStorage()
	if err != nil {
		log.Fatal("Failed to initialize storage because", err)
	}

	notifChannel, err := slack.NewNotification(config.Slack.Token)
	if err != nil {
		log.Fatal("Failed to initialize notification because", err)
	}

	flow := usecase.New(incidentStorage, notifChannel)

	handlers := handler.New(flow)

	router := chi.NewRouter()
	router.Get("/ping", ping)

	// Collection of datadog webhooks
	datadogWebhook := router.Group(nil)
	datadogWebhook.Route("/webhook/datadog", func(r chi.Router) {
		r.Post("/reply-in-thread", handlers.DdogReplyInThread)
	})

	grafanaWebhook := router.Group(nil)
	grafanaWebhook.Route("/webhook/grafana", func(r chi.Router) {
		//r.Post("/reply-in-thread", hOncall.GrafanaReplyInThread)
	})

	// Collection of datadog webhooks
	newrelicWebhook := router.Group(nil)
	newrelicWebhook.Route("/webhook/newrelic", func(r chi.Router) {
		//r.Post("/reply-in-thread", hOncall.NewRelicReplyInThread)
	})

	srv := http.Server{
		Addr:         config.Server.Port,
		ReadTimeout:  config.Server.ReadTimeout * time.Second,
		WriteTimeout: config.Server.WriteTimeout * time.Second,
		Handler:      router,
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Info("Server is up and running. Ready to receive request at", config.Server.Port)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Shutting down service...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully shut down")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

//-----------[ Helper ]-----------------

var logLevel = map[string]log.Level{
	"error": log.ErrorLevel,
	"warn":  log.WarnLevel,
	"info":  log.InfoLevel,
	"debug": log.DebugLevel,
}

func initLog(level, format string) error {
	lvl := log.ErrorLevel
	if value, ok := logLevel[level]; ok {
		lvl = value
	}
	log.SetLevel(lvl)

	if format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	return nil
}

func readConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config %s because %s", path, err)
	}

	config := &Config{}
	err = json.Unmarshal(content, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config %s because %s", path, err)
	}

	return config, nil
}
