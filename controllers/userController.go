package controllers

import (
	"database/sql"
	"encoding/json"
	"learn-api/config"
	"learn-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Param id query int false "User ID"
// @Param Authorization header string true "Authorization Bearer Token"
// @Success 200 {array} models.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบ Header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ตัดคำว่า "Bearer " ออกและเก็บ token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// ตรวจสอบ JWT Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_SECRET, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ตรวจสอบ query parameter id
	userIDParam := r.URL.Query().Get("id")
	query := "SELECT id, username, email FROM user"
	var args []interface{}

	// เพิ่มเงื่อนไขใน query หากมีการระบุ id
	if userIDParam != "" {
		query += " WHERE id = ?"
		id, err := strconv.Atoi(userIDParam)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}
		args = append(args, id)
	}

	// รัน query
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "Invalid username or password"
// @Failure 500 {string} string "Internal server error"
// @Router /api/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	var user models.User
	err := config.DB.QueryRow("SELECT id, username, password FROM user WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token หมดอายุใน 24 ชั่วโมง
	})

	tokenString, err := token.SignedString(config.JWT_SECRET)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
