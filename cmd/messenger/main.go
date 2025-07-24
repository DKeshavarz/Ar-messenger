package main

import (
	"log"
	"net/http"

	"github.com/DKeshavarz/Ar-messenger/internal/handlers"
	"github.com/DKeshavarz/Ar-messenger/internal/services"
	"github.com/gorilla/mux"
)

func main(){
	
	svc := services.NewRoomService(nil) // Replace nil with an actual MessageRepository implementation
    handler := handlers.NewWebSocketHandler(svc)

    router := mux.NewRouter()
    router.HandleFunc("/{chatName}/username", handler.HandleWebSocket)


    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}