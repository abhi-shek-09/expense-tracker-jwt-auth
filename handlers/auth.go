package handlers

import (
    "encoding/json"
    "expense-tracker/database"
    "expense-tracker/models"
    "expense-tracker/utils"
    "net/http"
    "strings"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
    err = database.DB.QueryRow(query, user.Username, user.Email, hashedPassword).Scan(&user.ID)
    if err != nil {
        if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
            http.Error(w, "Email already exists", http.StatusConflict)
        } else {
            http.Error(w, "Database error", http.StatusInternalServerError)
        }
        return
    }

    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(struct {
        Token string `json:"token"`
    }{Token: token})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
    var creds models.User
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    var user models.User
    query := "SELECT id, password FROM users WHERE email=$1"
    err = database.DB.QueryRow(query, creds.Email).Scan(&user.ID, &user.Password)
    if err != nil {
        http.Error(w, "Didn't find email", http.StatusUnauthorized)
        return
    }

    if !utils.CheckPassword(user.Password, creds.Password){
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(struct {
        Token string `json:"token"`
    }{Token: token})
}
