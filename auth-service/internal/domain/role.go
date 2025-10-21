package domain

type Role string

const (
	RoleGuest   Role = "guest"
	RoleUser    Role = "user"
	RoleManager Role = "manager"
	RoleAdmin   Role = "admin"
)