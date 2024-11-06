package controllers

import (
	"database/sql"
	"encoding/json"
	"learn-api/config"
	"learn-api/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Get all users
// @Description Get a list of all users with pagination
// @Tags users
// @Produce json
// @Param id query int false "User ID"
// @Param page query int false "Page number"
// @Param size query int false "Number of items per page"
// @Param Authorization header string true "Authorization Bearer Token"
// @Success 200 {array} models.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/users [get]
// GetUsers handles fetching users with optional ID filtering, pagination, and JWT authorization
func GetUsers(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_SECRET, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDParam := r.URL.Query().Get("id")
	baseQuery := "SELECT id, username, email, first_name, last_name, user_type_id, branch_id, profile_image_path, createDate, updateDate, status_id, role_id, nickname, phone_number, edit_by, u_index, u_code, permission FROM user"
	countQuery := "SELECT COUNT(*) FROM user"
	var args []interface{}
	var conditions []string

	if userIDParam != "" {
		conditions = append(conditions, "id = ?")
		id, err := strconv.Atoi(userIDParam)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}
		args = append(args, id)
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	pageParam := r.URL.Query().Get("page")
	sizeParam := r.URL.Query().Get("size")
	page := 1
	size := 10

	if pageParam != "" {
		page, err = strconv.Atoi(pageParam)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}

	if sizeParam != "" {
		size, err = strconv.Atoi(sizeParam)
		if err != nil || size < 1 {
			http.Error(w, "Invalid size parameter", http.StatusBadRequest)
			return
		}
	}

	offset := (page - 1) * size
	baseQuery += " LIMIT ? OFFSET ?"
	args = append(args, size, offset)

	// คำนวณจำนวนรายการทั้งหมด
	var totalElements int
	row := config.DB.QueryRow(countQuery, args[:len(args)-2]...)
	if err := row.Scan(&totalElements); err != nil {
		http.Error(w, "Error counting users", http.StatusInternalServerError)
		return
	}

	totalPages := (totalElements + size - 1) / size

	rows, err := config.DB.Query(baseQuery, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createDate, updateDate []byte

		if err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName,
			&user.UserTypeID, &user.BranchID, &user.ProfileImagePath, &createDate, &updateDate,
			&user.StatusID, &user.RoleID, &user.Nickname, &user.PhoneNumber, &user.EditBy,
			&user.UIndex, &user.UCode, &user.Permission,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.CreateDate, err = time.Parse("2006-01-02 15:04:05", string(createDate))
		if err != nil {
			log.Printf("Failed to parse createDate: %v", err)
			http.Error(w, "Error parsing createDate", http.StatusInternalServerError)
			return
		}

		user.UpdateDate, err = time.Parse("2006-01-02 15:04:05", string(updateDate))
		if err != nil {
			log.Printf("Failed to parse updateDate: %v", err)
			http.Error(w, "Error parsing updateDate", http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	response := map[string]interface{}{
		"totalPages":    totalPages,
		"totalElements": totalElements,
		"currentPage":   page,
		"content":       users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
	err := config.DB.QueryRow("SELECT id, username, password, email FROM user WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Email)

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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"token": tokenString,
	})
}
