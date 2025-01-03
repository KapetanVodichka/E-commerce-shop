package repositories

import (
	"database/sql"
	"e-commerce-shop/internal/model/good"
	"fmt"
)

type GoodRepository struct {
	db *sql.DB
}

func NewGoodRepository(db *sql.DB) *GoodRepository {
	return &GoodRepository{db: db}
}

func (r *GoodRepository) CreateGood(g *good.Good) error {
	if g.Discount < 0 {
		g.Discount = 0
	}

	g.TotalPrice = g.BasePrice * (1 - float64(g.Discount)/100)

	query := `
		INSERT INTO goods 
			(title, description, base_price, colour, size, count, discount, total_price)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.db.QueryRow(query,
		g.Title,
		g.Description,
		g.BasePrice,
		g.Colour,
		g.Size,
		g.Count,
		g.Discount,
		g.TotalPrice,
	).Scan(&g.ID)
	if err != nil {
		return fmt.Errorf("ошибка при создании товара: %w", err)
	}

	return nil
}

func (r *GoodRepository) GetGoodByID(id int) (*good.Good, error) {
	query := `
		SELECT id, title, description, base_price, colour, size, count, discount, total_price
		FROM goods
		WHERE id = $1
	`

	var g good.Good
	err := r.db.QueryRow(query, id).Scan(
		&g.ID,
		&g.Title,
		&g.Description,
		&g.BasePrice,
		&g.Colour,
		&g.Size,
		&g.Count,
		&g.Discount,
		&g.TotalPrice,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при попытке получить нужный товар: %w", err)
	}

	return &g, nil
}

func (r *GoodRepository) GetAllGoods() ([]good.Good, error) {
	query := `
		SELECT id, title, description, base_price, colour, size, count, discount, total_price
		FROM goods
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка товаров: %w", err)
	}
	defer rows.Close()

	var goods []good.Good
	for rows.Next() {
		var g good.Good
		if err := rows.Scan(
			&g.ID,
			&g.Title,
			&g.Description,
			&g.BasePrice,
			&g.Colour,
			&g.Size,
			&g.Count,
			&g.Discount,
			&g.TotalPrice,
		); err != nil {
			return nil, fmt.Errorf("ошибка при получении списка товаров: %w", err)
		}
		goods = append(goods, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при получении списка товаров: %w", err)
	}

	return goods, nil
}

func (r *GoodRepository) UpdateGood(id int, g *good.Good) error {
	if g.Discount < 0 {
		g.Discount = 0
	}

	g.TotalPrice = g.BasePrice * (1 - float64(g.Discount)/100)

	query := `
		UPDATE goods
		SET title = $1, description = $2, base_price = $3, colour = $4, size = $5, count = $6, discount = $7, total_price = $8
		WHERE id = $9
	`

	res, err := r.db.Exec(query,
		g.Title,
		g.Description,
		g.BasePrice,
		g.Colour,
		g.Size,
		g.Count,
		g.Discount,
		g.TotalPrice,
		id,
	)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении товара: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при обновлении товара: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ошибка при обновлении товара: товар с id %d не найден", id)
	}

	return nil
}

func (r *GoodRepository) DeleteGood(id int) error {
	query := `DELETE FROM goods WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении товара: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при удалении товара: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ошибка при удалении товара: товар с id %d не найден", id)
	}

	return nil
}
