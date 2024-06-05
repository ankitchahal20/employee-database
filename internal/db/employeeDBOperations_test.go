package db

import (
	"assignment/internal/constants"
	"assignment/internal/models"
	"assignment/internal/utils"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateEmployee_Success(t *testing.T) {
	utils.InitLogClient()

	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a new instance of the postgres struct with the mock database
	p := postgres{db: mockDB}

	// Create a test context and request
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Set up the expected SQL query and result
	var salary float64 = 50000.0
	employee := models.Employee{Name: "John Doe", Position: "Engineer", Salary: &salary}
	mock.ExpectQuery(`INSERT INTO employees \(name, position, salary\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(employee.Name, employee.Position, employee.Salary).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Call the CreateEmployee function
	employeeID, employeeErr := p.CreateEmployee(ctx, employee)

	// Assert that there's no error
	assert.Nil(t, employeeErr)

	// Assert the returned employee ID
	assert.Equal(t, "1", employeeID)

	// Assert that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestCreateEmployee_Error(t *testing.T) {
	utils.InitLogClient()

	// Create a new mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a new instance of the postgres struct with the mock database
	p := postgres{db: mockDB}

	// Set up the expected SQL query to return an error
	var salary float64 = 50000.0
	employee := models.Employee{Name: "John Doe", Position: "Engineer", Salary: &salary}
	mock.ExpectQuery(`INSERT INTO employees \(name, position, salary\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(employee.Name, employee.Position, employee.Salary).
		WillReturnError(errors.New("database error"))

	// Create a test context and request
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Call the CreateEmployee function
	_, employeeErr := p.CreateEmployee(ctx, employee)

	// Assert that there's an error
	assert.NotNil(t, employeeErr)

	// Assert that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateEmployee_Success(t *testing.T) {
	// Create a new mock database
	utils.InitLogClient()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a new instance of the postgres struct with the mock database
	p := postgres{db: mockDB}
	var salary float64 = 50000.0
	// Set up the expected SQL query and result
	employee := models.Employee{
		ID:       "1",
		Name:     "Updated Name",
		Position: "Updated Position",
		Salary:   &salary, // Assuming salary is updated
	}
	mock.ExpectExec(`UPDATE employees SET name=\$1, position=\$2, salary=\$3, last_updated_at=\$4 WHERE id=\$5`).
		WithArgs(employee.Name, employee.Position, *employee.Salary, sqlmock.AnyArg(), employee.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Indicating one row affected

		// Create a test context and request
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Call the UpdateEmployee function
	updatedEmployee, employeeErr := p.UpdateEmployee(ctx, employee)

	// Assert that there's no error
	assert.Nil(t, employeeErr)

	// Assert the returned employee
	assert.Equal(t, employee, updatedEmployee)

	// Assert that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateEmployee_NoFields(t *testing.T) {
	// Create a new mock database
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a new instance of the postgres struct with the mock database
	p := postgres{db: mockDB}

	// Create an employee with no fields to update
	employee := models.Employee{
		ID: "1",
	}

	// Set up a test context with a transaction ID
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Call the UpdateEmployee function
	_, employeeErr := p.UpdateEmployee(ctx, employee)

	// Assert that there's an error
	assert.NotNil(t, employeeErr)

	// Assert the error code and message
	assert.Equal(t, http.StatusBadRequest, employeeErr.Code)
	assert.Equal(t, "No fields to update", employeeErr.Message)
}

func TestUpdateEmployee_Error(t *testing.T) {
	// Create a new mock database
	utils.InitLogClient()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a new instance of the postgres struct with the mock database
	p := postgres{db: mockDB}

	var salary float64 = 50000.0
	// Set up the expected SQL query to return an error
	employee := models.Employee{
		ID:       "1",
		Name:     "Updated Name",
		Position: "Updated Position",
		Salary:   &salary,
	}
	mock.ExpectExec(`UPDATE employees SET name=\$1, position=\$2, salary=\$3, last_updated_at=\$4 WHERE id=\$5`).
		WithArgs(employee.Name, employee.Position, *employee.Salary, sqlmock.AnyArg(), employee.ID).
		WillReturnError(errors.New("database error"))

		// Set up a test context with a transaction ID
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Call the UpdateEmployee function
	_, employeeErr := p.UpdateEmployee(ctx, employee)

	// Assert that there's an error
	assert.NotNil(t, employeeErr)

	// Assert the error message
	assert.Equal(t, http.StatusInternalServerError, employeeErr.Code)
	assert.Equal(t, "Unable to update employee record", employeeErr.Message)
}

func TestGetEmployeeByID_Success(t *testing.T) {
	// Initialize logger and create a mock database
	utils.InitLogClient()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Create a new instance of the postgres struct with the mock database
	p := postgres{db: mockDB}

	// Define test data
	var salary float64 = 50000.0
	employeeID := "1"
	expectedEmployee := models.Employee{
		ID:            employeeID,
		Name:          "John Doe",
		Position:      "Engineer",
		Salary:        &salary,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	}

	// Set up the expected SQL query and result
	mock.ExpectQuery(`SELECT id, name, position, salary, created_at, last_updated_at FROM employees WHERE id=\$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "position", "salary", "created_at", "last_updated_at"}).
			AddRow(expectedEmployee.ID, expectedEmployee.Name, expectedEmployee.Position, expectedEmployee.Salary, expectedEmployee.CreatedAt, expectedEmployee.LastUpdatedAt))

	// Create a test context with a transaction ID
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.TransactionID: []string{"test-transaction-id"},
			},
		},
	}

	// Call the GetEmployeeByID function
	employee, employeeErr := p.GetEmployeeByID(ctx, employeeID)

	// Assert that there's no error
	assert.Nil(t, employeeErr)

	// Assert the returned employee
	assert.Equal(t, expectedEmployee, employee)

	// Assert that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
