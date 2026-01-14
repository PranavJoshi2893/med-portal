package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PranavJoshi2893/med-portal/internal/model"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user model.CreateUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user model.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
}

// func GetByID(w http.ResponseWriter, r *http.Request) {

// }

// func GetAll(w http.ResponseWriter, r *http.Request) {

// }

// func UpdateByID(w http.ResponseWriter, r *http.Request) {

// }

// func DeleteByID(w http.ResponseWriter, r *http.Request) {

// }
