package models

import "time"

// Employee struct defines the structure of an employee record
type Employee struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Position      string    `json:"position"`
	Salary        float64   `json:"salary"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}
