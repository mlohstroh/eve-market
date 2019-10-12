package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// Server is the struct for holding handles for
// everything like databases
type Server struct {
	ctx   context.Context
	oauth *oauth2.Config
}

func main() {
	godotenv.Load()

	profileFunction("loadSDE", func() {
		err := loadSDE()
		if err != nil {
			panic(err)
		}
	})

	profileFunction("getPrices", func() {
		err := getPrices()
		if err != nil {
			fmt.Print(err)
		}
	})

	scheduler := NewScheduler(1 * time.Hour)
	scheduler.Schedule("FetchPrices", getPrices, 5*time.Minute)
	go scheduler.Run()

	server := newServer()

	router := gin.Default()
	router.GET("/items/", getItems)
	router.GET("/items/:id", getItem)

	router.GET("/oauth/begin", server.oauthBegin)
	router.GET("/oauth/callback", server.oauthCallback)

	router.Run(":3000")
}

func newServer() *Server {
	clientID := os.Getenv("ESI_CLIENT_ID")
	if len(clientID) == 0 {
		panic("ESI_CLIENT_ID is not set")
	}

	secret := os.Getenv("ESI_SECRET_KEY")
	if len(secret) == 0 {
		panic("ESI_SECRET_KEY is not set")
	}

	callbackURL := os.Getenv("ESI_CALLBACK_URL")
	if len(callbackURL) == 0 {
		panic("ESI_CALLBACK_URL is not set")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		RedirectURL:  callbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/oauth/authorize",
			TokenURL: "https://login.eveonline.com/oauth/token",
		},
		Scopes: []string{"esi-markets.structure_markets.v1", "esi-markets.read_character_orders.v1", "esi-industry.read_character_jobs.v1", "esi-industry.read_character_mining.v1", "esi-contracts.read_character_contracts.v1", "esi-characters.read_blueprints.v1", "esi-assets.read_assets.v1", "esi-skills.read_skills.v1", "esi-skills.read_skillqueue.v1", "esi-ui.open_window.v1", "esi-wallet.read_character_wallet.v1"},
	}

	return &Server{
		ctx:   context.Background(),
		oauth: config,
	}
}
