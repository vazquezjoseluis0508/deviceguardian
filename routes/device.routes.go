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
	r.HandleFunc("/devices/{id}", UpdateDeviceHandler).Methods("PUT")
	r.HandleFunc("/devices/{id}", DeleteDeviceHandler).Methods("DELETE")
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

// En routes/devices.routes.go o un archivo similar

func UpdateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID del dispositivo de la URL
	vars := mux.Vars(r)
	deviceID := vars["id"]

	// Verificar que el dispositivo pertenece al usuario logueado
	userID, _ := auth.ExtractUserIDFromRequest(r)

	var updateDevice models.Device
	if err := json.NewDecoder(r.Body).Decode(&updateDevice); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, nil, "Invalid request body")
		return
	}

	var device models.Device
	if err := db.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
		utils.RespondJSON(w, http.StatusNotFound, nil, "Device not found")
		return
	}

	// Actualizar el dispositivo
	db.DB.Model(&device).Updates(updateDevice)
	utils.RespondJSON(w, http.StatusOK, device, nil)
}

// En routes/devices.routes.go o un archivo similar

func DeleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["id"]

	userID, _ := auth.ExtractUserIDFromRequest(r)

	// Intentar eliminar el dispositivo
	if err := db.DB.Where("id = ? AND user_id = ?", deviceID, userID).Delete(&models.Device{}).Error; err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, "Error deleting device")
		return
	}

	utils.RespondJSON(w, http.StatusOK, nil, "Device deleted successfully")
}
