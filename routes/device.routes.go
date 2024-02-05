package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vazquezjoseluis0508/go-gorm-api/auth"
	"github.com/vazquezjoseluis0508/go-gorm-api/db"
	"github.com/vazquezjoseluis0508/go-gorm-api/models"
	"github.com/vazquezjoseluis0508/go-gorm-api/utils"
)

type DeviceRequest struct {
	Name             string `json:"name" validate:"required"`
	Type             string `json:"type" validate:"required"`
	Description      string `json:"description"`
	Location         string `json:"location"`
	InstallationDate string `json:"installation_date"`
	Status           string `json:"status" validate:"required"`
}

func RegisterDeviceRoutes(r *mux.Router) {
	r.HandleFunc("/devices", CreateDeviceHandler).Methods("POST")
	r.HandleFunc("/devices", ListDevicesHandler).Methods("GET")
	// Agrega aquí las rutas para actualizar y eliminar dispositivos
}

// Handler para crear un nuevo dispositivo
func CreateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var req DeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, nil, map[string]string{"error": err.Error()})
		return
	}

	errors := utils.ValidateStruct(req)
	if len(errors) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, nil, map[string]interface{}{"errors": errors})
		return
	}

	userID, err := auth.ExtractUserIDFromRequest(r)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, nil, "Unauthorized or invalid user "+err.Error())
		return
	}

	device := models.Device{
		Name:             req.Name,
		Type:             req.Type,
		Description:      req.Description,
		Location:         req.Location,
		InstallationDate: req.InstallationDate,
		Status:           req.Status,
		UserID:           userID,
	}

	if result := db.DB.Create(&device); result.Error != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, "Error creating device"+result.Error.Error())
		return
	}

	utils.RespondJSON(w, http.StatusCreated, device, nil)
}

// Handler para listar dispositivos de un usuario
func ListDevicesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserIDFromRequest(r)
	if err != nil {
		// Manejar error, por ejemplo, devolviendo una respuesta de error
		http.Error(w, "Unauthorized or invalid user", http.StatusUnauthorized)
		return
	}
	var devices []models.Device

	if result := db.DB.Where("user_id = ?", userID).Find(&devices); result.Error != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, "Error listing devices")
		return
	}

	utils.RespondJSON(w, http.StatusOK, devices, nil)
}

// Implementa aquí los handlers para actualizar y eliminar dispositivos
