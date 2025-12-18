package app

import (
	"PrepaGo/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func AddEntityAToEntityB(w http.ResponseWriter, r *http.Request) {
	entityBIdStr := r.PathValue("entityBId")
	entityAIdStr := r.PathValue("entityAId")

	entityBId, err := strconv.Atoi(entityBIdStr)
	if err != nil {
		http.Error(w, "Invalid EntityB ID", http.StatusBadRequest)
		return
	}

	entityAId, err := strconv.Atoi(entityAIdStr)
	if err != nil {
		http.Error(w, "Invalid EntityA ID", http.StatusBadRequest)
		return
	}

	entityBExists, err := db.CheckEntityBExists(entityBId)
	if err != nil {
		http.Error(w, "Error checking EntityB", http.StatusInternalServerError)
		return
	}
	if !entityBExists {
		http.Error(w, "EntityB not found", http.StatusNotFound)
		return
	}

	entityAExists, err := db.CheckEntityAExists(entityAId)
	if err != nil {
		http.Error(w, "Error checking EntityA", http.StatusInternalServerError)
		return
	}
	if !entityAExists {
		http.Error(w, "EntityA not found", http.StatusNotFound)
		return
	}

	relationExists, err := db.CheckRelationExists(entityBId, entityAId)
	if err != nil {
		http.Error(w, "Error checking relation", http.StatusInternalServerError)
		return
	}
	if relationExists {
		http.Error(w, "Relation already exists", http.StatusConflict)
		return
	}

	entityB, err := db.GetEntityBById(entityBId)
	if err != nil {
		http.Error(w, "Error fetching EntityB", http.StatusInternalServerError)
		return
	}

	if !entityB.BoolEx {
		count, err := db.CountEntityAForEntityB(entityBId)
		if err != nil {
			http.Error(w, "Error counting relations", http.StatusInternalServerError)
			return
		}
		if count >= 5 {
			http.Error(w, "Non-bool EntityB can only have 5 EntityA maximum", http.StatusForbidden)
			return
		}
	}

	err = db.AddEntityAToEntityB(entityBId, entityAId)
	if err != nil {
		http.Error(w, "Error creating relation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "EntityA added to EntityB successfully")
}

func RemoveEntityAFromEntityB(w http.ResponseWriter, r *http.Request) {
	entityBIdStr := r.PathValue("entityBId")
	entityAIdStr := r.PathValue("entityAId")

	entityBId, err := strconv.Atoi(entityBIdStr)
	if err != nil {
		http.Error(w, "Invalid EntityB ID", http.StatusBadRequest)
		return
	}

	entityAId, err := strconv.Atoi(entityAIdStr)
	if err != nil {
		http.Error(w, "Invalid EntityA ID", http.StatusBadRequest)
		return
	}

	err = db.RemoveEntityAFromEntityB(entityBId, entityAId)
	if err != nil {
		http.Error(w, "Relation not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "EntityA removed from EntityB successfully")
}

func GetAllEntityAForEntityB(w http.ResponseWriter, r *http.Request) {
	entityBIdStr := r.PathValue("entityBId")
	entityBId, err := strconv.Atoi(entityBIdStr)
	if err != nil {
		http.Error(w, "Invalid EntityB ID", http.StatusBadRequest)
		return
	}

	exists, err := db.CheckEntityBExists(entityBId)
	if err != nil {
		http.Error(w, "Error checking EntityB", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "EntityB not found", http.StatusNotFound)
		return
	}

	entities, err := db.GetAllEntityAForEntityB(entityBId)
	if err != nil {
		http.Error(w, "Error fetching entities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entities)
}