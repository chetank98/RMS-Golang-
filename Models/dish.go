package Models

type DishCreation struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type DishRequest struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Price        int    `json:"price" db:"price"`
	RestaurantID string `json:"restaurantId" db:"restaurant_id"`
}
