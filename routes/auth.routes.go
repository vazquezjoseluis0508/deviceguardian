// auth.routes.go
package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vazquezjoseluis0508/go-gorm-api/auth"
	"github.com/vazquezjoseluis0508/go-gorm-api/db"
	"github.com/vazquezjoseluis0508/go-gorm-api/models"
	"github.com/vazquezjoseluis0508/go-gorm-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func RegisterAuthRoutes(r *mux.Router) {
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, nil, map[string]string{"error": err.Error()})
		return
	}

	errors := utils.ValidateStruct(req)
	if len(errors) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, nil, map[string]interface{}{"errors": errors})
		return
	}

	// Validar que el email no esté ya en uso
	var user models.User
	result := db.DB.Where("email = ?", req.Email).First(&user)
	if result.Error == nil {
		// Si no hay error, significa que encontramos un usuario con ese email
		utils.RespondJSON(w, http.StatusBadRequest, nil, map[string]string{"error": "This email is already in use."})
		return
	}

	// Crear el usuario en la base de datos si el email no está en uso
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, map[string]string{"error": "Error hashing the password."})
		return
	}

	user = models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
	}
	result = db.DB.Create(&user)
	if result.Error != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, map[string]string{"error": "Error creating user: " + result.Error.Error()})
		return
	}

	// Ocultar la contraseña en la respuesta
	user.Password = ""
	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{"user": user}, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	var user models.User
	result := db.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, nil, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, nil, "Invalid email or password")
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, "Error generating token")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"token": token}, nil)
}
