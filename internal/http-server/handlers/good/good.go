package good

import (
	"database/sql"
	"e-commerce-shop/internal/model/good"
	"e-commerce-shop/internal/storage/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func CreateGood(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var g good.Good
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(g); err != nil {
			http.Error(w, "Ошибка валидации: "+err.Error(), http.StatusBadRequest)
			return
		}

		repo := repositories.NewGoodRepository(db)
		if err := repo.CreateGood(&g); err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code == "23505" { // unique_violation
					http.Error(w, "Товар уже существует", http.StatusConflict)
					return
				}
			}
			http.Error(w, "Не удалось создать товар: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := map[string]interface{}{
			"id":          g.ID,
			"total_price": g.TotalPrice,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func GetGoodDetail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		goodIDStr := chi.URLParam(r, "goodID")
		goodID, err := strconv.Atoi(goodIDStr)
		if err != nil {
			http.Error(w, "Неверный идентификатор товара", http.StatusBadRequest)
			return
		}

		repo := repositories.NewGoodRepository(db)
		good, err := repo.GetGoodByID(goodID)
		if err != nil {
			http.Error(w, "Товар не найден: "+err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(good)
	}
}

func GetGoodList(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := repositories.NewGoodRepository(db)
		goods, err := repo.GetAllGoods()
		if err != nil {
			http.Error(w, "Не удалось получить список товаров: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(goods)
	}
}

func ChangeGood(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		goodIDStr := chi.URLParam(r, "goodID")
		goodID, err := strconv.Atoi(goodIDStr)
		if err != nil {
			http.Error(w, "Неверный идентификатор товара", http.StatusBadRequest)
			return
		}

		var g good.Good
		if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		validate := validator.New()
		if err := validate.Struct(g); err != nil {
			http.Error(w, "Ошибка валидации: "+err.Error(), http.StatusBadRequest)
			return
		}

		repo := repositories.NewGoodRepository(db)
		if err := repo.UpdateGood(goodID, &g); err != nil {
			http.Error(w, "Не удалось обновить товар: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"message":     fmt.Sprintf("Товар с ID %d обновлён", goodID),
			"total_price": g.TotalPrice,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteGood(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		goodIDStr := chi.URLParam(r, "goodID")
		goodID, err := strconv.Atoi(goodIDStr)
		if err != nil {
			http.Error(w, "Неверный идентификатор товара", http.StatusBadRequest)
			return
		}

		repo := repositories.NewGoodRepository(db)
		if err := repo.DeleteGood(goodID); err != nil {
			http.Error(w, "Не удалось удалить товар: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"message": fmt.Sprintf("Товар с ID %d удалён", goodID),
		}
		json.NewEncoder(w).Encode(response)
	}
}
