package main

import (
	g "auth_service/grpc"
	h "auth_service/http"
	"auth_service/http/api"
	"auth_service/storage"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	storage.ConnectToDB()

	defer storage.DB.Close()

	go g.GrpcListen()

	router := mux.NewRouter()

	router.HandleFunc("/is-auth", api.IsAuthorizedHandler)
	router.HandleFunc("/login", api.LoginHandler)
	router.HandleFunc("/register", api.RegisterHandler)
	router.HandleFunc("/users/{username}", api.UserHandler)

	http.Handle("/", router)

	listen := os.Getenv("AUTH_SERVICE_LISTEN_HTTP_IP") + ":" + os.Getenv("AUTH_SERVICE_LISTEN_HTTP_PORT")
	server := h.NewServer(listen)

	fmt.Println("http server is listening on port " + listen)
	log.Fatal(server.ListenAndServe())
}
