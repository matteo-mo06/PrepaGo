package app

import (
	"encoding/json"
	"exam_api/db"
	"exam_api/models"
	"fmt"
	"net/http"
	"strconv"
)

func GetAllEntityB(w http.ResponseWriter, r *http.Request) {
	entities, err := db.GetAllEntityB()
	if err != nil {
		http.Error(w, "Error fetching entities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entities)
}

func GetEntityBById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	entity, err := db.GetEntityBById(id)
	if err != nil {
		http.Error(w, "Entity not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func CreateEntityB(w http.ResponseWriter, r *http.Request) {
	var entity models.EntityB
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var errMsgs []string

	if len(entity.Name) < 2 || len(entity.Name) > 100 {
		errMsgs = append(errMsgs, "Name must be between 2 and 100 characters")
	}

	if len(entity.MaxVarcharEx) > 255 {
		errMsgs = append(errMsgs, "MaxVarcharEx must not exceed 255 characters")
	}

	if len(errMsgs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string][]string{"errors": errMsgs})
		return
	}

	if entity.Id != 0 {
		exists, err := db.CheckEntityBExists(entity.Id)
		if err != nil {
			http.Error(w, "Error checking ID", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "ID already exists", http.StatusConflict)
			return
		}
	}

	lastId, err := db.CreateEntityB(entity)
	if err != nil {
		http.Error(w, "Error creating entity", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": lastId})
}

func UpdateEntityBBoolEx(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	exists, err := db.CheckEntityBExists(id)
	if err != nil {
		http.Error(w, "Error checking entity", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Entity not found", http.StatusNotFound)
		return
	}

	var requestBody struct {
		BoolEx bool `json:"bool_ex"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = db.UpdateEntityBBoolEx(id, requestBody.BoolEx)
	if err != nil {
		http.Error(w, "Error updating bool_ex", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "BoolEx updated successfully")
}

func DeleteEntityB(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = db.DeleteEntityB(id)
	if err != nil {
		http.Error(w, "Entity not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Entity deleted successfully")
}