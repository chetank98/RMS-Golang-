package Models

type Role string

const (
	RoleAdmin    Role = "Admin"
	RoleSubAdmin Role = "Sub-admin"
	RoleUser     Role = "User"
)

type UserRequest struct {
	Name     string        `json:"name" db:"name"`
	Email    string        `json:"email" db:"email"`
	Password string        `json:"password" db:"password"`
	Address  []UserAddress `json:"address" db:"address"`
}

type User struct {
	ID      string    `json:"id" db:"id"`
	Name    string    `json:"name" db:"name"`
	Email   string    `json:"email" db:"email"`
	Address []Address `json:"address" db:"address"`
	Role    Role      `json:"role" db:"role"`
}

type UserAddress struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Address struct {
	ID        string  `json:"id" db:"id"`
	Address   string  `json:"address" db:"address"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
	UserID    string  `json:"userId" db:"user_id"`
}

type LoginRequest struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type LoginDetail struct {
	ID       string `json:"email" db:"id"`
	Password string `json:"password" db:"password"`
	Role     Role   `json:"role" db:"role"`
}

type UserCtx struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Role      Role   `json:"role"`
}

type DistanceRequest struct {
	UserAddressID       string `json:"userAddressId" validate:"required"`
	RestaurantAddressID string `json:"restaurantAddressId" validate:"required"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
