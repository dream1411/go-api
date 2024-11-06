package routes

import (
	userController "learn-api/controllers"

	"learn-api/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/login", userController.Login).Methods("POST")

	// สร้างกลุ่ม API ที่ต้องการตรวจสอบสิทธิ์ด้วย JWT
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JWTAuthMiddleware) // ใช้ Middleware JWTAuthMiddleware กับกลุ่มนี้
	api.HandleFunc("/users", userController.GetUsers).Methods("GET")
	return r
}
