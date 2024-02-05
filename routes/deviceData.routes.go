package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vazquezjoseluis0508/go-gorm-api/db"
	"github.com/vazquezjoseluis0508/go-gorm-api/models"
	"github.com/vazquezjoseluis0508/go-gorm-api/utils"
	"github.com/vazquezjoseluis0508/go-gorm-api/websocket"
)

type DeviceDataRequest struct {
	DeviceID      uint    `json:"device_id" validate:"required"`
	Timestamp     string  `json:"timestamp" validate:"required"`
	Temperature   float64 `json:"temperature" validate:"required"`
	BatteryStatus string  `json:"battery_status" validate:"required"`
	// Añade aquí otros sensores o datos relevantes
}

// RegisterDeviceDataRoutes registra las rutas para el manejo de datos de dispositivos
func RegisterDeviceDataRoutes(r *mux.Router) {
	r.HandleFunc("/device-data", PostDeviceDataHandler).Methods("POST")
}

// PostDeviceDataHandler maneja la recepción de datos de dispositivo y envía actualizaciones por WebSocket
func PostDeviceDataHandler(w http.ResponseWriter, r *http.Request) {
	var req DeviceDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, nil, map[string]string{"error": err.Error()})
		return
	}

	errors := utils.ValidateStruct(req)
	if len(errors) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, nil, map[string]interface{}{"errors": errors})
		return
	}

	// Verificar si el dispositivo existe
	var device models.Device
	if err := db.DB.First(&device, req.DeviceID).Error; err != nil {
		utils.RespondJSON(w, http.StatusNotFound, nil, map[string]string{"error": "Device not found"})
		return
	}

	// Crear el DeviceData con la información proporcionada
	deviceData := models.DeviceData{
		DeviceID:      req.DeviceID,
		Timestamp:     utils.ParseTime(req.Timestamp),
		Temperature:   req.Temperature,
		BatteryStatus: req.BatteryStatus,
		// Añade aquí otros sensores o datos relevantes
	}

	// Guardar los datos del dispositivo en la base de datos
	if result := db.DB.Create(&deviceData); result.Error != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, nil, map[string]string{"error": result.Error.Error()})
		return
	}

	// Enviar los datos del dispositivo a todos los clientes conectados por WebSocket
	websocket.Broadcast(deviceData)

	utils.RespondJSON(w, http.StatusCreated, deviceData, nil)
}
