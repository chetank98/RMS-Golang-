package Models

type SubAdminRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gte=6,lte=15"`
}

type SubAdmin struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Role      Role   `json:"role" db:"role"`
	CreatedBy string `json:"createdBy" db:"created_by"`
}
