package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bishal05das/auth-service/internal/database"
	"github.com/bishal05das/auth-service/internal/models"
	"github.com/bishal05das/auth-service/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	db *database.Database
	jwtSecret []byte
	//Add token expiation configuration
	tokenExpiration time.Duration
}

func NewAuthHandler(db *database.Database, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		db: db,
		jwtSecret: jwtSecret,
		tokenExpiration: 24 * time.Hour,  //default 24 hour expiration
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request){
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var user models.UserRegister
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
			"details": err.Error(),
	    })
		return
	}

	var exists bool
	err := h.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)",user.Email).Scan(&exists)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Database error"})
		return
	}
	if exists {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "email already registered"})
		return
	}
	//hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "password processing failed"})
		return
	}
	tx,err := h.db.DB.Begin()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Transaction start failed"})
		return
	}
	var id int
	err = tx.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", user.Email, hashedPassword).Scan(&id)

	if err != nil {
		tx.Rollback()
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "User creation failed"})
		return
	}

	if err = tx.Commit(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Transaction commit failed"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"message": "User registered successfully", "user_id": id,})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request){
	var login models.UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&login); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error" : "Invalid request body",
			"details": err.Error(),
		})
	}
	//get user from database
	var user models.User
	err := h.db.DB.QueryRow("SELECT id, email, password_hash FROM users WHERE email=$1", login.Email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Login process failed"})
		return
	}

	//verify password
	if !utils.CheckPasswordHash(login.Password, user.PasswordHash) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid Password"})
		return
	}
    //generate jwt token
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"iat": now.Unix(),
		"exp": now.Add(h.tokenExpiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Token generation failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token": tokenString,
		"expires_in": h.tokenExpiration.Seconds(),
		"token_type": "Bearer",
	})


}

func writeJSON(w http.ResponseWriter, status int, data interface{}){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}


