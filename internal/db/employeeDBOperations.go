package db

import (
	employeeerror "assignment/internal/errors"
	"assignment/internal/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateEmployee function
func (p postgres) CreateEmployee(ctx *gin.Context, employee models.Employee) (string, *employeeerror.EmployeeError) {
	txid := ctx.Request.Header.Get("TransactionID")

	query := `INSERT INTO employees (name, position, salary) VALUES ($1, $2, $3) RETURNING id`
	var employeeID int

	err := p.db.QueryRowContext(ctx, query, employee.Name, employee.Position, employee.Salary).Scan(&employeeID)
	if err != nil {
		fmt.Printf("error while running insert query, txid: %v\n", txid)
		return "", &employeeerror.EmployeeError{
			Trace:   txid,
			Code:    http.StatusInternalServerError,
			Message: "unable to add employee",
		}
	}

	id := strconv.Itoa(employeeID)
	fmt.Printf("successfully added employee entry in db, txid: %v\n", txid)
	return id, nil
}



func (p postgres) DeleteEmployee(ctx *gin.Context, notes models.Employee) *employeeerror.EmployeeError {
	return nil
}

func (p postgres) GetEmployeeByID(ctx *gin.Context) ([]models.Employee, *employeeerror.EmployeeError) {
	return []models.Employee{}, nil
}
func (p postgres) UpdateEmployee(ctx *gin.Context) (string, *employeeerror.EmployeeError) {
	return "", nil
}