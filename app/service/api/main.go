package main

import (
	v1 "bbs-game/app/service/api/v1"
	"bbs-game/domain/account"
	"bbs-game/domain/battle"
	"bbs-game/domain/creature"
	"bbs-game/domain/item"
	"bbs-game/domain/player"
	"bbs-game/gateways/persistence"
	"bbs-game/middlewares"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	pgAddr := "postgres://ps_user:ps_password@localhost:7002/bbs-game?sslmode=disable"
	db, err := persistence.NewDB(pgAddr)
	if err != nil {
		log.Fatal(err)
	}

	itemService := item.NewService(db)
	itemAPI := v1.NewItemAPI(itemService)

	playerService := player.NewService(db)
	playerAPI := v1.NewPlayerAPI(playerService)

	creatureService := creature.NewService(db)
	creatureAPI := v1.NewCreatureAPI(creatureService)

	battleService := battle.NewService(db)
	battleAPI := v1.NewBattleAPI(battleService)

	accountService := account.NewService(db)
	accountAPI := v1.NewAccountAPI(accountService)

	motdAPI := v1.NewMotdAPI()

	credentials := handlers.AllowCredentials()
	headers := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})

	r := mux.NewRouter()
	r.Use(middlewares.Authentication)

	// Public Routes
	r.HandleFunc("/api/v1/motd", motdAPI.Get).Methods("GET")
	r.HandleFunc("/api/v1/accounts", accountAPI.Create).Methods("POST")
	r.HandleFunc("/api/v1/accounts/login", accountAPI.Login).Methods("POST")

	// Private Routes
	r.HandleFunc("/api/v1/battles", battleAPI.Create).Methods("POST")
	r.HandleFunc("/api/v1/battles/{id}", battleAPI.Get).Methods("GET")
	r.HandleFunc("/api/v1/battles/{id}/actions", battleAPI.CreateAction).Methods("POST")
	r.HandleFunc("/api/v1/battles/{id}/next_round", battleAPI.GenerateNextRound).Methods("GET")
	r.HandleFunc("/api/v1/battles/{id}/actions", battleAPI.ListActionsByBattle).Methods("GET")
	r.HandleFunc("/api/v1/items", itemAPI.Create).Methods("POST")
	r.HandleFunc("/api/v1/items/{id}", itemAPI.Get).Methods("GET")
	r.HandleFunc("/api/v1/players", playerAPI.Create).Methods("POST")
	r.HandleFunc("/api/v1/players/me", playerAPI.GetPlayerByToken).Methods("GET")
	r.HandleFunc("/api/v1/players/{id:[0-9]+}", playerAPI.Get).Methods("GET")
	r.HandleFunc("/api/v1/creatures", creatureAPI.Create).Methods("POST")
	r.HandleFunc("/api/v1/creatures/{id}", creatureAPI.Get).Methods("GET")

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", handlers.CORS(credentials, methods, origins, headers)(r)); err != nil {
		log.Fatal(err)
	}
}
