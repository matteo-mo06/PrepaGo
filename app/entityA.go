package app

import (
	"PrepaGo/db"
	"PrepaGo/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetAllEntityA(w http.ResponseWriter, r *http.Request) {
	entities, err := db.GetAllEntityA()
	if err != nil {
		http.Error(w, "Error fetching entities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entities)
}

func GetEntityAById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	entity, err := db.GetEntityAById(id)
	if err != nil {
		http.Error(w, "Entity not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func CreateEntityA(w http.ResponseWriter, r *http.Request) {
	var entity models.EntityA
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var errMsgs []string

	if len(entity.Name) < 4 || len(entity.Name) > 100 {
		errMsgs = append(errMsgs, "Name must be between 4 and 100 characters")
	}

	if entity.DecimalEx < 0 {
		errMsgs = append(errMsgs, "Decimal must be positive")
	}

	if len(errMsgs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string][]string{"errors": errMsgs})
		return
	}

	if entity.Id != 0 {
		exists, err := db.CheckEntityAExists(entity.Id)
		if err != nil {
			http.Error(w, "Error checking ID", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "ID already exists", http.StatusConflict)
			return
		}
	}

	lastId, err := db.CreateEntityA(entity)
	if err != nil {
		http.Error(w, "Error creating entity", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": lastId})
}

func UpdateEntityA(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	exists, err := db.CheckEntityAExists(id)
	if err != nil {
		http.Error(w, "Error checking entity", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Entity not found", http.StatusNotFound)
		return
	}

	var entity models.EntityA
	err = json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var errMsgs []string

	if len(entity.Name) < 4 || len(entity.Name) > 100 {
		errMsgs = append(errMsgs, "Name must be between 4 and 100 characters")
	}

	if entity.DecimalEx < 0 {
		errMsgs = append(errMsgs, "Decimal must be positive")
	}

	if len(errMsgs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string][]string{"errors": errMsgs})
		return
	}

	err = db.UpdateEntityA(id, entity)
	if err != nil {
		http.Error(w, "Error updating entity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Entity updated successfully")
}

func DeleteEntityA(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = db.DeleteEntityA(id)
	if err != nil {
		http.Error(w, "Entity not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Entity deleted successfully")
}