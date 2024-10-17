package Models

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleSubAdmin Role = "sub-admin"
	RoleUser     Role = "user"
)

type User struct {
	Id         string `json:"userId" db:"id"`
	Name       string `json:"userName" db:"name"`
	Email      string `json:"userEmail" db:"email"`
	Password   string `json:"password" db:"password"`
	Created_by string `json:"createdBy" db:"created_by"`
}

type Login struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Role     Role   `db:"role"`
}

type UserCtx struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Role      Role   `json:"role"`
}
