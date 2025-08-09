package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/DKeshavarz/Ar-messenger/internal/config"
	"github.com/DKeshavarz/Ar-messenger/internal/handlers"
	"github.com/DKeshavarz/Ar-messenger/internal/repositories"
	"github.com/DKeshavarz/Ar-messenger/internal/services"
	"github.com/gorilla/mux"
)

func main() {
	brokersEnv := config.GetEnvValue("REDPANDA_BROKERS")
	if brokersEnv == "" {
		log.Fatal("REDPANDA_BROKERS env not set")
	}
	brokers := strings.Split(brokersEnv, ",")
    
	repoRedpanda, err := repositories.NewRedpandaMessageRepository(brokers)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create redpanda repo: %w", err))
	}
	defer repoRedpanda.Close()

	svc := services.NewRoomService(repoRedpanda)
	handler := handlers.NewWebSocketHandler(svc)

	router := mux.NewRouter()
	router.HandleFunc("/{chatName}/username", handler.HandleWebSocket)
	
	http.Handle("/", router)
	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "web/index.html")
	// }).Methods("GET")
	
	log.Println("Server starting on :" + config.GetEnvValue("SERVER_PORT"))
	log.Fatal(http.ListenAndServe(":"+config.GetEnvValue("SERVER_PORT"), nil))
}