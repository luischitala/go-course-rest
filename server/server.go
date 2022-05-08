package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	database "github.com/luischitala/go-server/database"
	repository "github.com/luischitala/go-server/repository"
	websocket "github.com/luischitala/go-server/websocket"
	"github.com/rs/cors"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

//Will handle the serverss
type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

//We need to return the config
func (b *Broker) Config() *Config {
	return b.config
}

//Function to implement the hub
func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	//Review the configuration to assure that there are no empty fields
	if config.Port == "" {
		return nil, errors.New("Port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("Secret is required")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("Database is required")
	}
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}

	return broker, nil
}

//Method to allow the broker execute
func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	//How to handle the cors
	handler := cors.Default().Handler(b.router)
	//Database repository
	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	//Initialize the hub ()
	go b.hub.Run()
	repository.SetRepository(repo)
	log.Println("Starting server on port", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
