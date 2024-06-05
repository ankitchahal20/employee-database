package db

import (
	employeeerror "assignment/internal/errors"
	"assignment/internal/models"
	"database/sql"
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

func (p postgres) DeleteEmployee(ctx *gin.Context, employeeId string) *employeeerror.EmployeeError {
	txid := ctx.Request.Header.Get("TransactionID")
	
	// Convert employeeId to integer and handle any errors
	empId, _ := strconv.Atoi(employeeId)

	// SQL query to delete employee by ID
	query := `DELETE FROM employees WHERE id=$1`

	// Execute the query
	if _, err := p.db.ExecContext(ctx, query, empId); err != nil {
		fmt.Println("Error executing delete query, empId:", empId, "error:", err)
		return &employeeerror.EmployeeError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to delete employee record",
			Trace:   txid,
		}
	}

	fmt.Printf("Successfully deleted employee entry from db, txid: %v\n", txid)
	return nil
}


func (p postgres) GetEmployeeByID(ctx *gin.Context, employeeId string) (models.Employee, *employeeerror.EmployeeError) {
	txid := ctx.Request.Header.Get("TransactionID")
	fmt.Println("employeeId:", employeeId)

	// Convert employeeId to integer and handle any errors
	empId, _ := strconv.Atoi(employeeId)

	// SQL query to get employee by ID
	query := `SELECT id, name, position, salary, created_at, last_updated_at FROM employees WHERE id=$1`

	// Prepare to scan the result into an Employee struct
	employee := &models.Employee{}
	err := p.db.QueryRowContext(ctx, query, empId).Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary, &employee.CreatedAt, &employee.LastUpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Employee{}, &employeeerror.EmployeeError{
				Code:    http.StatusNotFound,
				Message: "Employee not found",
				Trace:   txid,
			}
		}
		fmt.Println("Error executing query, empId:", empId, "error:", err)
		return models.Employee{}, &employeeerror.EmployeeError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to retrieve employee record",
			Trace:   txid,
		}
	}

	fmt.Printf("Successfully retrieved employee entry from db, txid: %v\n", txid)
	return *employee, nil

}
func (p postgres) UpdateEmployee(ctx *gin.Context) (string, *employeeerror.EmployeeError) {
	return "", nil
}