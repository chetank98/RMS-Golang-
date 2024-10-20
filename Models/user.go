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

type UserAddress struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
