package router

import (
	"mosanda-vacation-attendace/handlers"
	"net/http"
)

func NewRouter(handlers *handlers.Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /api/login", handlers.Login)
	router.Handle("GET /api/me", handlers.RequireSession(http.HandlerFunc(handlers.GetCurrentUser)))
	router.Handle("GET /api/vacation-batches", handlers.RequireSession(http.HandlerFunc(handlers.GetVacationBatches)))
	router.Handle("POST /api/vacation-batches", handlers.RequireSession(http.HandlerFunc(handlers.CreateVacationBatch)))
	router.Handle("GET /api/vacation-batches/{id}", handlers.RequireSession(http.HandlerFunc(handlers.GetVacationBatch)))
	router.Handle("PUT /api/vacation-batches/{id}", handlers.RequireSession(http.HandlerFunc(handlers.UpdateVacationBatch)))
	router.Handle("DELETE /api/vacation-batches/{id}", handlers.RequireSession(http.HandlerFunc(handlers.DeleteVacationBatch)))
	router.Handle("GET /api/attendances", handlers.RequireSession(http.HandlerFunc(handlers.GetAttendances)))
	router.Handle("POST /api/attendances", handlers.RequireSession(http.HandlerFunc(handlers.CreateAttendance)))
	router.Handle("GET /api/attendances/{id}", handlers.RequireSession(http.HandlerFunc(handlers.GetAttendance)))
	router.Handle("PUT /api/attendances/{id}", handlers.RequireSession(http.HandlerFunc(handlers.UpdateAttendance)))
	router.Handle("DELETE /api/attendances/{id}", handlers.RequireSession(http.HandlerFunc(handlers.DeleteAttendance)))

	return router
}
