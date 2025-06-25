package main

import (
	"log"
	"net/http"

	"github.com/bishal05das/auth-service/internal/config"
	"github.com/bishal05das/auth-service/internal/database"
	"github.com/bishal05das/auth-service/internal/handlers"
)

func main() {
    
    cfg,err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:",err)
	}
	//initialize database
	db,err := database.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect database",err)
	}
	defer db.DB.Close()

	router := http.NewServeMux()
	//auth handler
	authHandler := handlers.NewAuthHandler(db, []byte(cfg.JWT.Secretkey))

	//routing
	router.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	router.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s\n", serverAddr)
	server := &http.Server{
		Addr: serverAddr,
		Handler: router,
		ReadTimeout: cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,

	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:",err)
	}
}
