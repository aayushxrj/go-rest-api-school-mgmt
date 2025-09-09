package models

import (
	"database/sql/driver"
	"errors"
)

// NullString is used here only for Swagger, but now supports DB scanning.
type NullString struct {
	String string `json:"string,omitempty"`
	Valid  bool   `json:"valid,omitempty"`
}

// Scan implements the sql.Scanner interface
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String = ""
		ns.Valid = false
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan NullString: incompatible type")
	}
	ns.String = string(b)
	ns.Valid = true
	return nil
}

// Value implements the driver.Valuer interface
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

type Exec struct {
	ID                   int        `json:"id,omitempty" db:"id,omitempty"`
	FirstName            string     `json:"first_name,omitempty" db:"first_name,omitempty"`
	LastName             string     `json:"last_name,omitempty" db:"last_name,omitempty"`
	Email                string     `json:"email,omitempty" db:"email,omitempty"`
	Username             string     `json:"username,omitempty" db:"username,omitempty"`
	Password             string     `json:"password,omitempty" db:"password,omitempty"`
	PasswordChangedAt    NullString `json:"password_changed_at,omitempty" db:"password_changed_at,omitempty"`
	UserCreatedAt        NullString `json:"user_created_at,omitempty" db:"user_created_at,omitempty"`
	PasswordResetToken   NullString `json:"password_reset_token,omitempty" db:"password_reset_token,omitempty"`
	PasswordTokenExpires NullString `json:"password_token_expires,omitempty" db:"password_token_expires,omitempty"`
	InactiveStatus       bool       `json:"inactive_status,omitempty" db:"inactive_status,omitempty"`
	Role                 string     `json:"role,omitempty" db:"role,omitempty"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword string `json:"new_password"`
}

type UpdatePasswordResponse struct {
	Token string `json:"token"`
	PasswordUpdated bool `json:"password_updated"`
}