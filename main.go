package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vazquezjoseluis0508/go-gorm-api/db"
	"github.com/vazquezjoseluis0508/go-gorm-api/models"
	"github.com/vazquezjoseluis0508/go-gorm-api/routes"
)

func main() {
	db.DBConnection()

	db.DB.AutoMigrate(&models.User{}, &models.Device{}, &models.DeviceData{}, &models.Alert{})

	r := mux.NewRouter()

	routes.RegisterRoutes(r)

	fmt.Println("Server running on port 3000")

	http.ListenAndServe(":3000", r)

}
