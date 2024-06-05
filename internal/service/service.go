package service

import (
	"assignment/internal/constants"
	"assignment/internal/db"
	employeeerror "assignment/internal/errors"
	"assignment/internal/models"
	"assignment/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

		txid := ctx.Request.Header.Get(constants.TransactionID)
		utils.Logger.Info(fmt.Sprintf("received request for employee creation, txid : %v", txid))
		var notes models.Employee
		if err := ctx.ShouldBindBodyWith(&notes, binding.JSON); err == nil {
			utils.Logger.Info(fmt.Sprintf("user request for employee creation is unmarshalled successfully, txid : %v", txid))
			employeeID, err := employeeClient.createEmployee(ctx, notes)
			if err != nil {
				utils.RespondWithError(ctx, err.Code, err.Message)
				return
			}
			ctx.JSON(http.StatusOK, map[string]string{
				"employee_id": employeeID,
			})
			ctx.Writer.WriteHeader(http.StatusOK)

		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"Unable to marshal the request body": err.Error()})
		}
	}
}

func (service *EmployeeService) createEmployee(ctx *gin.Context, employee models.Employee) (string, *employeeerror.EmployeeError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)

	utils.Logger.Info(fmt.Sprintf("calling db layer for employee creation, txid : %v", txid))
	employeeID, err := service.repo.CreateEmployee(ctx, employee)
	if err != nil {
		return "", err
	}
	return employeeID, nil
}

// Deletes an employee from the database or store by ID
func DeleteEmployee() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		employeeId := ctx.Param("id")
		err := employeeClient.deleteEmployee(ctx, employeeId)
		if err != nil {
			utils.RespondWithError(ctx, err.Code, err.Message)
			return
		}

		utils.Logger.Info(fmt.Sprintf("user has successfully deleted a note, txid : %v", txid))
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

func (service *EmployeeService) deleteEmployee(ctx *gin.Context, employeeId string) *employeeerror.EmployeeError {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	
	utils.Logger.Info(fmt.Sprintf("calling db layer for to check if employee Id exists, txid : %v", txid))
	_, err := service.repo.GetEmployeeByID(ctx, employeeId)
	if err != nil {
		return err
	}

	utils.Logger.Info(fmt.Sprintf("calling db layer for employee deletion, txid : %v", txid))
	err = service.repo.DeleteEmployee(ctx, employeeId)
	if err != nil {
		return err
	}
	return nil
}

// Retrieves an employee from the database or store by ID
func GetEmployeeByID() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		txid := ctx.Request.Header.Get(constants.TransactionID)
		employeeId := ctx.Param("id")
		utils.Logger.Info(fmt.Sprintf("calling service layer for to get the employee details, txid : %v", txid))
		employeeDetails, err := employeeClient.getEmployeeByID(ctx, employeeId)
		if err != nil {
			utils.RespondWithError(ctx, err.Code, err.Message)
			return
		}
		ctx.JSON(http.StatusOK, map[string]string{
			"employee_id": fmt.Sprintf("%v", employeeDetails.ID),
			"employee_name": employeeDetails.Name,
			"position": employeeDetails.Position,	
		})
	}
}

func (service *EmployeeService) getEmployeeByID(ctx *gin.Context, employeeId string) (models.Employee, *employeeerror.EmployeeError) {
	txid := ctx.Request.Header.Get(constants.TransactionID)
	
	utils.Logger.Info(fmt.Sprintf("calling db layer for to get employee details, txid : %v", txid))

	employeeDetails, err := service.repo.GetEmployeeByID(ctx, employeeId)
	if err != nil {
		return models.Employee{}, err
	}
	return employeeDetails, nil
}

// Updates the details of an existing employee
func UpdateEmployee() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
	}
}
