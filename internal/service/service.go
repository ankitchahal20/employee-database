package service

import (
	"assignment/internal/db"

	"github.com/gin-gonic/gin"
)

var (
	employeeClient *EmployeeService
)

type EmployeeService struct {
	repo db.EmployeeDBService
}

func NewEmployeeService(conn db.EmployeeDBService) *EmployeeService {
	employeeClient = &EmployeeService{
		repo: conn,
	}
	return employeeClient
}

// Adds a new employee to the database
func CreateEmployee() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
	}
}

// Deletes an employee from the database or store by ID
func DeleteEmployee() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
	}
}

// Retrieves an employee from the database or store by ID
func GetEmployeeByID() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
	}
}

// Updates the details of an existing employee
func UpdateEmployee() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
	}
}
