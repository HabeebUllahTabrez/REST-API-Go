package models

type User struct {
	Name        string `json:"name,omitempty" validate:"required"`
	DOB         string `json:"dob,omitempty" validate:"required"`
	Address     string `json:"address,omitempty" validate:"required"`
	Description string `json:"description,omitempty" validate:"required"`
	CreatedAt   string `json:"createdAt,omitempty"`
}
