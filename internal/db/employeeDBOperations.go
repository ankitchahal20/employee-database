package db

import (
	"assignment/internal/constants"
	employeeerror "assignment/internal/errors"
	"assignment/internal/models"
	"assignment/internal/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateEmployee function
func (p postgres) CreateEmployee(ctx *gin.Context, employee models.Employee) (string, *employeeerror.EmployeeError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

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
	utils.Logger.Info(fmt.Sprintf("successfully added employee entry in db, txid: %v\n", txid))
	return id, nil
}

func (p postgres) DeleteEmployee(ctx *gin.Context, employeeId string) *employeeerror.EmployeeError {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	// Convert employeeId to integer and handle any errors
	empId, _ := strconv.Atoi(employeeId)
	fmt.Println("empId : ", empId)
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

	utils.Logger.Info(fmt.Sprintf("Successfully deleted employee entry from db, txid: %v\n", txid))
	return nil
}

func (p postgres) GetEmployeeByID(ctx *gin.Context, employeeId string) (models.Employee, *employeeerror.EmployeeError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	fmt.Println("employeeId:", employeeId)

	// Convert employeeId to integer and handle any errors
	empId, err := strconv.Atoi(employeeId)
	if err != nil {
		fmt.Println("Error converting employee ID to integer:", err)
		return models.Employee{}, &employeeerror.EmployeeError{
			Code:    http.StatusBadRequest,
			Message: "Invalid employee ID",
			Trace:   txid,
		}
	}

	// SQL query to get employee by ID
	query := `SELECT id, name, position, salary, created_at, last_updated_at FROM employees WHERE id=$1`

	// Prepare to scan the result into an Employee struct
	employee := &models.Employee{}
	err = p.db.QueryRowContext(ctx, query, empId).Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary, &employee.CreatedAt, &employee.LastUpdatedAt)
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

	utils.Logger.Info(fmt.Sprintf("Successfully retrieved employee entry from db, txid: %v\n", txid))
	return *employee, nil
}

func (p postgres) UpdateEmployee(ctx *gin.Context, employee models.Employee) (models.Employee, *employeeerror.EmployeeError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	// Build the dynamic update query
	var fields []string
	var args []interface{}
	argID := 1

	if employee.Name != "" {
		fields = append(fields, fmt.Sprintf("name=$%d", argID))
		args = append(args, employee.Name)
		argID++
	}
	if employee.Position != "" {
		fields = append(fields, fmt.Sprintf("position=$%d", argID))
		args = append(args, employee.Position)
		argID++
	}
	if employee.Salary != nil {
		fields = append(fields, fmt.Sprintf("salary=$%d", argID))
		args = append(args, *employee.Salary)
		argID++
	}

	// If no fields to update, return an error
	if len(fields) == 0 {
		return models.Employee{}, &employeeerror.EmployeeError{
			Code:    http.StatusBadRequest,
			Message: "No fields to update",
			Trace:   txid,
		}
	}

	// Add the last_updated_at field if there are other fields being updated
	fields = append(fields, fmt.Sprintf("last_updated_at=$%d", argID))
	args = append(args, time.Now())
	argID++

	// Add the ID to the arguments
	args = append(args, employee.ID)

	query := fmt.Sprintf("UPDATE employees SET %s WHERE id=$%d", strings.Join(fields, ", "), argID)
	res, err := p.db.Exec(query, args...)
	if err != nil {
		fmt.Println("Error executing update query:", err)
		return models.Employee{}, &employeeerror.EmployeeError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to update employee record",
			Trace:   txid,
		}
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error checking rows affected:", err)
		return models.Employee{}, &employeeerror.EmployeeError{
			Code:    http.StatusInternalServerError,
			Message: "Unable to verify update operation",
			Trace:   txid,
		}
	}
	if rowsAffected == 0 {
		return models.Employee{}, &employeeerror.EmployeeError{
			Code:    http.StatusNotFound,
			Message: "Employee not found",
			Trace:   txid,
		}
	}

	utils.Logger.Info(fmt.Sprintf("Successfully updated employee entry in db, txid: %v\n", txid))
	return employee, nil
}

func (p postgres) ListEmployee(ctx *gin.Context, page int, pageSize int) ([]models.Employee, *employeeerror.EmployeeError) {
    txid := ctx.Request.Header.Get(constants.TransactionID)

    // Calculate the offset based on the page number and page size
    offset := (page - 1) * pageSize

    // SQL query to list employee records with pagination
    query := `SELECT id, name, position, salary, created_at, last_updated_at 
               FROM employees 
               ORDER BY id 
               LIMIT $1 OFFSET $2`

    // Execute the query with the specified page size and offset
    rows, err := p.db.QueryContext(ctx, query, pageSize, offset)
    if err != nil {
        fmt.Println("Error executing query:", err)
        return nil, &employeeerror.EmployeeError{
            Code:    http.StatusInternalServerError,
            Message: "Unable to retrieve employee records",
            Trace:   txid,
        }
    }
    defer rows.Close()

    // Iterate over the rows and scan the results into Employee structs
    var employees []models.Employee
    for rows.Next() {
        var employee models.Employee
        if err := rows.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary, &employee.CreatedAt, &employee.LastUpdatedAt); err != nil {
            fmt.Println("Error scanning row:", err)
            return nil, &employeeerror.EmployeeError{
                Code:    http.StatusInternalServerError,
                Message: "Error processing employee records",
                Trace:   txid,
            }
        }
        employees = append(employees, employee)
    }

    // Check for any errors encountered during iteration
    if err := rows.Err(); err != nil {
        fmt.Println("Error iterating over rows:", err)
        return nil, &employeeerror.EmployeeError{
            Code:    http.StatusInternalServerError,
            Message: "Error processing employee records",
            Trace:   txid,
        }
    }

    utils.Logger.Info(fmt.Sprintf("Successfully retrieved employee records from db (page %d, pageSize %d), txid: %v\n", page, pageSize, txid))
    return employees, nil
}