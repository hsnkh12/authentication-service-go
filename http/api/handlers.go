package api

import (
	j "auth_service/jwt"
	"auth_service/storage/models"
	"auth_service/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {

	// Check header
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only DELETE method is allowed for '/users'"))
		return
	}
	//

	// Get user
	vars := mux.Vars(r)
	username := vars["username"]

	resUser, err := models.GetUserByUsername(username)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not found"))
		return
	}
	//

	// Parse auth token
	authToken := r.Header.Get("Authorization")
	tkn, err := j.ParseToken(authToken)

	if err != nil || tkn == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := tkn.Claims.(*j.Claim)

	if claims.User_id != resUser.Id {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//

	// Send request to chat service to delete user chat history
	chat_service_url := "http://" + os.Getenv("CHAT_SERVICE_LISTEN_HTTP_IP") + ":" + os.Getenv("CHAT_SERVICE_LISTEN_HTTP_PORT") + "/user-del-history/" + resUser.Id.String()
	request, err := http.NewRequest("GET", chat_service_url, nil)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("failed to create 'user-del-history' request: %v", err)
		return
	}

	request.Header.Add("Service-Communication-Token", os.Getenv("SERVICE_COM_TOKEN"))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("failed to send 'user-del-history' request to chat service: %v", err)
		return
	}
	//

	w.WriteHeader(response.StatusCode)

}

func IsAuthorizedHandler(w http.ResponseWriter, r *http.Request) {

	// Check header
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only GET method is allowed for '/is-auth'"))
		return
	}

	// Check JWT token if exists
	authToken := r.Header.Get("Authorization")
	tkn, err := j.ParseToken(authToken)

	if err != nil || tkn == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//

	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// Check header
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only POST method is allowed for '/login'"))
		return
	}

	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content-Type must be 'application/json'"))
		return
	}
	//

	// Check JWT token if exists
	authToken := r.Header.Get("Authorization")
	tkn, err := j.ParseToken(authToken)

	if err == nil && tkn != nil {

		if tkn.Valid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//

	// Read body
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Could not read content of body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	//

	// Parse body to JSON
	user := models.User{}

	err = json.Unmarshal(body, &user)

	if err != nil {
		log.Println("Body could not be parsed to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//

	// Get the user
	resUser, err := models.GetUserByUsername(user.Username)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not found"))
		return
	}
	//

	// Check password
	err = utils.ComparePassword(resUser.Password, user.Password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
		return
	}
	//

	// Sign new JWT token
	expirationTime := time.Now().Add(time.Minute * 30)

	claim := &j.Claim{User_id: resUser.Id, StandardClaims: jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
	}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"token":"` + tokenString + `"}`))
	w.WriteHeader(http.StatusAccepted)

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	// Check header
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only POST method is allowed for '/register'"))
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content-Type must be 'application/json'"))
		return
	}
	//

	// Read body
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Could not read content of body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	//

	// Parse body to JSON
	user := models.User{}

	err = json.Unmarshal(body, &user)

	if err != nil {
		log.Println("Body could not be parsed to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//

	// Store user to DB
	err = user.StoreToDB()

	if err != nil {

		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Username must be unique"))
			return
		}

		fmt.Println(user)
		log.Panic(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	//

	w.WriteHeader(http.StatusOK)

}
