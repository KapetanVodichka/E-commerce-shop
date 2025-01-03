package good

type Good struct {
	ID          int     `json:"id" db:"id"`
	Title       string  `json:"title" db:"title" validate:"required,min=3,max=255"`
	Description string  `json:"description" db:"description"`
	BasePrice   float64 `json:"base_price" db:"base_price" validate:"required,gte=0"`
	Colour      string  `json:"colour" db:"colour" validate:"max=50"`
	Size        string  `json:"size" db:"size" validate:"max=50"`
	Count       int     `json:"count" db:"count" validate:"required,gte=0"`
	Discount    int     `json:"discount" db:"discount" validate:"gte=0,lte=100"`
	TotalPrice  float64 `json:"total_price" db:"total_price"`
}
