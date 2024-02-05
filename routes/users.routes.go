package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vazquezjoseluis0508/go-gorm-api/db"
	"github.com/vazquezjoseluis0508/go-gorm-api/models"
)

func RegisterUserRoutes(s *mux.Router) {
	s.HandleFunc("/users", GetUsersHandler).Methods("GET")
	s.HandleFunc("/users/{id}", GetUserHandler).Methods("GET")
	s.HandleFunc("/users", CreateUserHandler).Methods("POST")
	s.HandleFunc("/users/{id}", UpdateUserHandler).Methods("PUT")
	s.HandleFunc("/users/{id}", DeleteUserHandler).Methods("DELETE")
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.DB.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if params["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID is required"))
		return
	}

	var user models.User
	response := db.DB.First(&user, params["id"])
	err := response.Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(&user)
	}

}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	createdUser := db.DB.Create(&user)
	err := createdUser.Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&user)
	}

}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	if params["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID is required"))
		return
	}

	response := db.DB.Delete(&models.User{}, params["id"])
	err := response.Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("User deleted"))
	}

}
