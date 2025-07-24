package main

import (
	"log"
	"net/http"

	"github.com/DKeshavarz/Ar-messenger/internal/config"
	"github.com/DKeshavarz/Ar-messenger/internal/handlers"
	"github.com/DKeshavarz/Ar-messenger/internal/services"
	"github.com/gorilla/mux"
)

func main(){
	svc := services.NewRoomService(nil) // Replace nil with an actual MessageRepository implementation
    handler := handlers.NewWebSocketHandler(svc)

    router := mux.NewRouter()
    router.HandleFunc("/{chatName}/username", handler.HandleWebSocket)
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "web/index.html")
    }).Methods("GET")
    router.HandleFunc("/{chatName}/username", handler.HandleWebSocket)

    http.Handle("/", router)
    http.Handle("/style.css", http.FileServer(http.Dir("web")))
    http.Handle("/main.js", http.FileServer(http.Dir("web")))

    log.Println("Server starting on :" + config.GetEnvValue("SERVER_PORT"))
    log.Fatal(http.ListenAndServe(":" + config.GetEnvValue("SERVER_PORT") , nil))
}