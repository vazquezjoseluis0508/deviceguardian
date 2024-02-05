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

type AlertRequest struct {
	DeviceID    uint   `json:"device_id" validate:"required"`
	AlertType   string `json:"alert_type" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func RegisterAlertRoutes(r *mux.Router) {
	r.HandleFunc("/alerts", CreateAlertHandler).Methods("POST")
	r.HandleFunc("/alerts", ListAlertsByUserHandler).Methods("GET")
}

func CreateAlertHandler(w http.ResponseWriter, r *http.Request) {
	var req AlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, nil, map[string]string{"error": err.Error()})
		return
	}

	userID, _ := auth.ExtractUserIDFromRequest(r)

	errors := utils.ValidateStruct(req)
	if len(errors) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, nil, map[string]interface{}{"errors": errors})
		return
	}

	alert := models.Alert{
		DeviceID:    req.DeviceID,
		UserID:      userID,
		AlertType:   req.AlertType,
		Description: req.Description,
	}

	if result := db.DB.Create(&alert); result.Error != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, map[string]string{"error": result.Error.Error()})
		return
	}

	utils.RespondJSON(w, http.StatusCreated, alert, nil)
}

func ListAlertsByUserHandler(w http.ResponseWriter, r *http.Request) {
	var alerts []models.Alert
	userID, _ := auth.ExtractUserIDFromRequest(r)

	// Aquí, podrías añadir lógica para filtrar las alertas por userID, DeviceID, etc.
	result := db.DB.Where("user_id = ?", userID).Find(&alerts)
	if result.Error != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, result.Error.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, alerts, nil)
}

// Handler para actualizar un dispositivo
func DeleteAlertHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alertID := vars["id"]

	var alert models.Alert
	result := db.DB.First(&alert, alertID)
	if result.Error != nil {
		utils.RespondJSON(w, http.StatusNotFound, nil, map[string]string{"error": "Alert not found"})
		return
	}

	if result := db.DB.Delete(&alert); result.Error != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, map[string]string{"error": result.Error.Error()})
		return
	}

	utils.RespondJSON(w, http.StatusOK, nil, nil)
}
