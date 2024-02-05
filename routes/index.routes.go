// index.routes.go
package routes

import (
	"github.com/gorilla/mux"
	"github.com/vazquezjoseluis0508/go-gorm-api/middleware"
)

func RegisterRoutes(r *mux.Router) {
	RegisterAuthRoutes(r)       // Registrar rutas de autenticación
	RegisterDeviceDataRoutes(r) // Registrar rutas de datos de dispositivos

	// Rutas protegidas
	s := r.PathPrefix("/api").Subrouter()
	s.Use(middleware.JWTAuthentication)
	RegisterUserRoutes(s)
	RegisterDeviceRoutes(s)

}
