package Models

type RestaurantsRequest struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Restaurant struct {
	ID        string  `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	Address   string  `json:"address" db:"address"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
	CreatedBy string  `json:"createdBy" db:"created_by"`
}
