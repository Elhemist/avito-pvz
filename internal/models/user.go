package models

import "github.com/google/uuid"

type Role string

const (
	RoleEmployee  Role = "employee"
	RoleModerator Role = "moderator"
	RoleClient    Role = "client"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Role         Role      `db:"role"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    string    `db:"created_at"`
}
