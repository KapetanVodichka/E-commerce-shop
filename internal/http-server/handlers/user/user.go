package user

import (
	"database/sql"
	"e-commerce-shop/internal/model/user"
	"e-commerce-shop/internal/storage/repositories"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req user.User
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			http.Error(w, "Ошибка валидации: "+err.Error(), http.StatusBadRequest)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка хеширования пароля: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Password = string(hashed)

		repo := repositories.NewUserRepository(db)
		if err := repo.CreateUser(&req); err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code == "23505" {
					http.Error(w, "Пользователь с таким именем уже существует", http.StatusConflict)
					return
				}
			}
			http.Error(w, "Ошибка создания пользователя: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resp := map[string]interface{}{
			"id": req.ID,
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type loginReq struct {
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		}
		var req loginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			http.Error(w, "Ошибка валидации: "+err.Error(), http.StatusBadRequest)
			return
		}

		repo := repositories.NewUserRepository(db)
		u, err := repo.GetUserByUsername(req.Username)
		if err != nil {
			http.Error(w, "Неправильное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
			http.Error(w, "Неправильное имя пользователя или пароль", http.StatusUnauthorized)
			return
		}

		token, err := createToken(u.ID, u.Username)
		if err != nil {
			http.Error(w, "Не удалось создать токен: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"token": token,
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func createToken(userID int, username string) (string, error) {
	key := secretKey
	if envKey := os.Getenv("JWT_SECRET"); envKey != "" {
		key = []byte(envKey)
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}
