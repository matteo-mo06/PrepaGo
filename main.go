package main

import (
	"PrepaGo/app"
	"PrepaGo/db"
	"fmt"
	"net/http"
)

func main() {
	db.Conn = db.NewDB()
	defer db.Conn.Close()

	http.HandleFunc("GET /{$}", healthCheck)

	http.HandleFunc("GET /entity-a/{$}", app.GetAllEntityA)
	http.HandleFunc("GET /entity-a/{id}", app.GetEntityAById)
	http.HandleFunc("POST /entity-a/{$}", app.CreateEntityA)
	http.HandleFunc("PUT /entity-a/{id}", app.UpdateEntityA)
	http.HandleFunc("DELETE /entity-a/{id}", app.DeleteEntityA)

	http.HandleFunc("GET /entity-b/{$}", app.GetAllEntityB)
	http.HandleFunc("GET /entity-b/{id}", app.GetEntityBById)
	http.HandleFunc("POST /entity-b/{$}", app.CreateEntityB)
	http.HandleFunc("PATCH /entity-b/{id}/bool", app.UpdateEntityBBoolEx)
	http.HandleFunc("DELETE /entity-b/{id}", app.DeleteEntityB)

	http.HandleFunc("POST /entity-b/{entityBId}/entity-a/{entityAId}", app.AddEntityAToEntityB)
	http.HandleFunc("DELETE /entity-b/{entityBId}/entity-a/{entityAId}", app.RemoveEntityAFromEntityB)
	http.HandleFunc("GET /entity-b/{entityBId}/entity-a/{$}", app.GetAllEntityAForEntityB)

	fmt.Println("Server running on http://localhost:8090")
	http.ListenAndServe(":8090", nil)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	err := db.Conn.Ping()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "API is running")
}