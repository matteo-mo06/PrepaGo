package app

import (
	"PrepaGo/db"
	"PrepaGo/models"
	"PrepaGo/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func GetAllLibraries(w http.ResponseWriter, r *http.Request) {
	libraries, err := db.GetAllLibraries()
	if err != nil {
		http.Error(w, "Error fetching libraries", http.StatusInternalServerError)
		return
	}

	var response []map[string]any
	for _, lib := range libraries {
		response = append(response, map[string]any{
			"id":            lib.Id,
			"owner_name":    lib.OwnerName,
			"is_premium":    lib.IsPremium,
			"creation_year": lib.CreationYear,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetLibraryById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	library, err := db.GetLibraryById(id)
	if err != nil {
		http.Error(w, "Library not found", http.StatusNotFound)
		return
	}

	response := map[string]any{
		"id":            library.Id,
		"owner_name":    library.OwnerName,
		"is_premium":    library.IsPremium,
		"creation_year": library.CreationYear,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateLibrary(w http.ResponseWriter, r *http.Request) {
	var library models.Library
	err := json.NewDecoder(r.Body).Decode(&library)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var errMsgs []string

	if len(library.OwnerName) < 2 || len(library.OwnerName) > 100 {
		errMsgs = append(errMsgs, "Owner name must be between 2 and 100 characters")
	}

	if len(library.OwnerPassword) < 5 || len(library.OwnerPassword) > 100 {
		errMsgs = append(errMsgs, "Owner password must be between 5 and 100 characters")
	}

	if len(errMsgs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string][]string{"errors": errMsgs})
		return
	}

	if library.Id != 0 {
		exists, err := db.CheckLibraryExists(library.Id)
		if err != nil {
			http.Error(w, "Error checking ID", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "ID already exists", http.StatusConflict)
			return
		}
	}

	lastId, err := db.CreateLibrary(library)
	if err != nil {
		http.Error(w, "Error creating library", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": lastId})
}

func UpdateLibraryPremium(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	username, err := utils.VerifyJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	library, err := db.GetLibraryById(id)
	if err != nil {
		http.Error(w, "Library not found", http.StatusNotFound)
		return
	}

	if library.OwnerName != username {
		http.Error(w, "Forbidden: You can only modify your own library", http.StatusForbidden)
		return
	}

	var requestBody struct {
		IsPremium bool `json:"is_premium"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = db.UpdateLibraryPremium(id, requestBody.IsPremium)
	if err != nil {
		http.Error(w, "Error updating premium status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Premium status updated successfully")
}

func DeleteLibrary(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = db.DeleteLibrary(id)
	if err != nil {
		http.Error(w, "Library not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Library deleted successfully")
}