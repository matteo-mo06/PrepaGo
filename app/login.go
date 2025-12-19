package app

import (
	"PrepaGo/db"
	"PrepaGo/models"
	"PrepaGo/utils"
	"encoding/json"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	existing, err := db.GetLibraryByOwnerName(creds.Username)
	if err != nil {
		http.Error(w, "Error requesting library by owner name", http.StatusInternalServerError)
		return
	}

	if existing == nil || creds.Username != existing.OwnerName || creds.Password != existing.OwnerPassword {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}